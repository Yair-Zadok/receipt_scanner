# Copyright 2024, Yair Zadok, All rights reserved.

import base64
import requests
import csv
import os
from PIL import Image
from pdf2image import convert_from_path
import pandas as pd


api_key = ""

# General information prompt
text_prompt = r"""IT IS EXTREMELY IMPORTANT THAT ALL INSTRUCTIONS BE FOLLOWED WITH PRECISION: 
Given the provided image, give the subtotal, total, tax, date, and tips if they are present 
(otherwise leave the fields as XXX) in EXACTLY the following format, 
IF THE TAX IS INLCUDED IN THE SUBTOTAL, SUBTRACT IT FROM THE SUBTOTAL, GIVE NO OTHER RESPONSE BUT THE FOLLOWING: 
subtotal: XXX, total: XXX, tax: XXX, tips: XXX date: YEAR/MONTH/DAY
"""

# Prompt for matching to expense category
accounts_list = r"[Advertising, Bank charges, Dues and Subscriptions, Insurance, Legal and professional fees, Meals and entertainment, Office expenses, Promotional, Rent or lease payments, Office Supplies, Repair and maintenance, Taxes and Licenses, Travel, Travel meals, Uncategorized Expense, Utilities]"
text_prompt_accounts = f"""IT IS EXTREMELY IMPORTANT THAT ALL INSTRUCTIONS BE FOLLOWED WITH PRECISION: 
Given the provided image, pick which category best suits the receipt from the following list of categories,
for example purchase of a monitor from Staples would go under Office expenses {accounts_list}. 
Pick exactly one element in the list and return ONLY that element and NOTHING ELSE, 
return it in exactly the same format, capitalization and spacing, 
for example if your chosen element is 'Utilities' return ONLY: Utilities
"""

# Prompt for matching to a supplier
supplier_list = "Costco, Bills Accounting, Office Depot, Jims Lawfirm"
text_prompt_supplier = f"""IT IS EXTREMELY IMPORTANT THAT ALL INSTRUCTIONS BE FOLLOWED WITH PRECISION: 
Given the provided image, Decide if the receipts Vendor is already included in the following list of exisiting vendors, 
and if not return the name of the receipt's new unincluded vendor, {supplier_list}, PICK EXACTLY ONE VENDOR, 
and return it in exactly the same format, capitalization, and spacing, for example if your chosen vendor is 
'Williams Fresh Cafe' return ONLY: Williams Fresh Cafe
"""


# Gives a base64 encoding of an image
def encode_image(image_path : str) -> str:
    with open(image_path, "rb") as image_file:
        return base64.b64encode(image_file.read()).decode('utf-8')


# Retrieves a string response according to a text prompt and image path
def get_receipt_string(image_path : str, text_prompt : str, api_key : str) -> str:
    base64_image = encode_image(image_path)

    headers = {
      "Content-Type": "application/json",
      "Authorization": f"Bearer {api_key}"
    }

    payload = {
      "model": "gpt-4-vision-preview",
      "messages": [
        {
          "role": "user",
          "content": [
            {
              "type": "text",
              "text": f"{text_prompt}"
            },
            {
              "type": "image_url",
              "image_url": {
                "url": f"data:image/jpeg;base64,{base64_image}"
              }
            }
          ]
        }
      ],
      "max_tokens": 100
    }
    
    response = requests.post("https://api.openai.com/v1/chat/completions", headers=headers, json=payload)
    response_data = response.json()
    completion = response_data['choices'][0]['message']['content']
    return completion


class ReceiptData:
    def __init__(self, subtotal, total, tax, tips, date):
        self.subtotal = subtotal
        self.total = total
        self.tax = tax
        self.tips = tips
        self.date = date


# Converts a string prompt response to a ReceiptData struct
def convert_string_to_receipt_data(data_string : str) -> ReceiptData:
    try:
        data_string = data_string.replace('$', '')
        data_string = data_string.lower()
        data_string = data_string.replace('xxx', '0')
        data_parts = data_string.split(', ')
        subtotal = data_parts[0].split(': ')[1]
        total = data_parts[1].split(': ')[1]
        tax = data_parts[2].split(': ')[1]
        tips = data_parts[3].split(': ')[1]
        date = data_parts[4].split(': ')[1]
        return ReceiptData(subtotal, total, tax, tips, date)
    except Exception as e:
        print("An error occurred:", str(e))
        return ReceiptData("FAILED", "FAILED", "FAILED", "FAILED", "FAILED")
    

# Converts a ReceiptData struct into CSV form
def receipt_data_to_csv(receipt_data : ReceiptData, account : str, supplier : str,
                                image_path : str, csv_data : list[list[str]]=[]) -> list[list[str]]:
    if not csv_data:
        csv_data = []
        csv_data.append(['BillNo', 'supplier', 'account', 'subtotal', 'tax', 'tips', 'total', 'date', 'receipt_path'])

    csv_data.append([len(csv_data), supplier, account, receipt_data.subtotal, receipt_data.tax, receipt_data.tips, receipt_data.total, receipt_data.date, os.path.basename(image_path)])

    return csv_data


# Saves a CSV
def save_csv(csv_data : list[list[str]], filename : str):
    with open(filename, 'w', newline='') as csvfile:
        writer = csv.writer(csvfile)
        writer.writerows(csv_data)


# Some receipts may have multiple PDF pages and this function appends them vertically
def concatenate_images_vertically(images):
    total_height = sum(image.size[1] for image in images)
    
    max_width = max(image.size[0] for image in images)
    
    concatenated_image = Image.new("RGB", (max_width, total_height), color="white")
    
    y_offset = 0
    for image in images:
        concatenated_image.paste(image, (0, y_offset))
        y_offset += image.size[1]
    
    return concatenated_image


def scan_folder(folder_path : str, poppler_path : str, api_key : str):
    file_list = os.listdir(folder_path)
    csv_data = []
    
    # Iterate over every receipt in the target receipt folder
    for file_name in file_list:
        if file_name.endswith('.pdf'):
            # Get PDF image into variable form
            pdf_path = os.path.join(folder_path, file_name)
            images = convert_from_path(pdf_path, dpi=450, poppler_path=poppler_path)

            concatenated_image = concatenate_images_vertically(images)

            # Convert to PNG image
            png_path = os.path.join(folder_path, f"{os.path.splitext(file_name)[0]}.png")
            concatenated_image.save(png_path, "PNG")

            os.remove(pdf_path)
            
            # Call for AI reponse on prompts
            receipt_string = get_receipt_string(png_path, text_prompt, api_key)
            receipt_data = convert_string_to_receipt_data(receipt_string)
            account = get_receipt_string(png_path, text_prompt_accounts, api_key)
            supplier = get_receipt_string(png_path, text_prompt_supplier, api_key)

            csv_data = receipt_data_to_csv(receipt_data, account, supplier, png_path, csv_data)
            
        elif file_name.endswith(('.jpg', '.jpeg', '.png')):
            image_path = os.path.join(folder_path, file_name)

            # Call for AI reponse on prompts
            receipt_string = get_receipt_string(image_path, text_prompt, api_key)
            receipt_data = convert_string_to_receipt_data(receipt_string)
            account = get_receipt_string(image_path, text_prompt_accounts, api_key)
            supplier = get_receipt_string(image_path, text_prompt_supplier, api_key)

            csv_data = receipt_data_to_csv(receipt_data, account, supplier, image_path, csv_data)
            
    return csv_data

current_dir = os.path.dirname(__file__)
folder_path = os.path.join(current_dir, "receipts")
output_path = os.path.join(current_dir, "data_intake.csv")

poppler_path = r"POPPLER_PATH\Library\bin"
csv_data = scan_folder(folder_path, poppler_path, text_prompt, api_key)
save_csv(csv_data, output_path)
