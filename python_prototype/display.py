# Copyright 2024, Yair Zadok, All rights reserved.

import csv
from datetime import datetime
import os
from tkinter import messagebox
import tkinter as tk
import customtkinter as ctk
from PIL import Image, ImageTk


# Section: Data intake from CSV
class ReceiptData:
    def __init__(self, supplier, account, subtotal, tax, tips, total, date, receipt_path):
        self.supplier = supplier
        self.account = account
        self.subtotal = subtotal
        self.tax = tax
        self.tips = tips
        self.total = total
        self.date = date
        self.receipt_path = receipt_path


# Reads a data intake CSV into a list of ReceiptData structs
def read_receipt_csv(file_name : str) -> list[ReceiptData]:
    receipt_data_list = []

    with open(file_name, 'r') as file:
        reader = csv.DictReader(file)
        for row in reader:
            supplier = row['supplier']
            account = row['account']
            subtotal = row['subtotal']
            tax = row['tax']
            tips = row['tips']
            total = row['total']
            date_str = row['date']
            receipt_path = row['receipt_path']
            try:
                date = datetime.strptime(date_str, '%Y/%m/%d').date()
                date = date.strftime('%Y/%m/%d')
            except:
                date = date_str
            
            receipt_data = ReceiptData(supplier, account, subtotal, tax, tips, total, date, receipt_path)
            receipt_data_list.append(receipt_data)

    return receipt_data_list


# Section: TKinter Display

# Saves CSV data
def save_csv(csv_data : str, filename : str):
    with open(filename, 'w', newline='') as csvfile:
        writer = csv.writer(csvfile)
        writer.writerows(csv_data)


# Adds a TKinter entry field
def add_entry(label_text : str):
    label = ctk.CTkLabel(left_frame, text=label_text, font=("Arial", 15))
    label.pack(side=tk.TOP, anchor="nw", padx=200, pady=2)
    entry = ctk.CTkEntry(left_frame, width=270, font=("Arial", 19))
    entry.pack(side=tk.TOP, anchor="nw", padx=200, pady=0)
    
    return label, entry


# Pre-populates the entry fields with data from the data_intake.csv
def populate_entries(receipt_data : ReceiptData):
    supplier_entry.delete(0, tk.END)
    supplier_entry.insert(tk.END, receipt_data.supplier)
    
    account_entry.delete(0, tk.END)
    account_entry.insert(tk.END, receipt_data.account)
    
    subtotal_entry.delete(0, tk.END)
    subtotal_entry.insert(tk.END, receipt_data.subtotal)
    
    tax_entry.delete(0, tk.END)
    tax_entry.insert(tk.END, receipt_data.tax)
    
    tips_entry.delete(0, tk.END)
    tips_entry.insert(tk.END, receipt_data.tips)
    
    date_entry.delete(0, tk.END)
    date_entry.insert(tk.END, receipt_data.date)


# Handles swapping to next receipt image
def swap_photo(photo_path : str):
    current_dir = os.path.dirname(__file__)
    photo_folder = os.path.join(current_dir, "receipts")
    photo_path = os.path.join(photo_folder, photo_path)
    
    right_frame.update()
    frame_width = right_frame.winfo_reqwidth()
    frame_height = right_frame.winfo_reqheight()
    photo = Image.open(photo_path)
    
    aspect_ratio = photo.width / photo.height
    
    if frame_width / aspect_ratio <= frame_height:
        width = frame_width
        height = width / aspect_ratio
    else:
        height = frame_height
        width = height * aspect_ratio
    
    photo = photo.resize((int(width), int(height)), Image.ANTIALIAS)

    tk_photo = ImageTk.PhotoImage(photo)
    
    image_label.configure(image=tk_photo)
    image_label.image = tk_photo 


# Converts entries in fields to CSV format
def convert_entries_to_csv(csv_data : list[list[str]]=[]) -> list[list[str]]:
    if not csv_data:
        csv_data = []
        csv_data.append(['BillNo', 'Supplier', 'BillDate', 'DueDate', 'Account', 'LineAmount', 'LineTaxCode', 'Line Tax Amount'])

    try:
        line_amount = float(subtotal_entry.get()) + float(tips_entry.get())
    except:
        line_amount = "FAILED"

    line_tax_amount = tax_entry.get()
    csv_data.append([len(csv_data), supplier_entry.get(), date_entry.get(), '2050/01/01', account_entry.get(), line_amount, 'HST ON', line_tax_amount])

    return csv_data


# Save data to output CSV
def save_data():
    try:
        global csv_data
        csv_data = convert_entries_to_csv(csv_data)
        supplier_entry.update()
        print(supplier_entry.get())
        current_dir = os.path.dirname(__file__)
        output_path = os.path.join(current_dir, "output.csv")
        save_csv(csv_data, output_path)

    except Exception as ed:
        messagebox.showerror("Error", f"An error occurred: {str(ed)}")


# Load next entry, handling receipt image swaps, data saving, and prep for next entry
def load_next(receipt_data_list : list[ReceiptData]):
    global global_index
    global_index += 1
    save_data()
 
    if global_index >= len(receipt_data_list):
        root.destroy()

    swap_photo(receipt_data_list[global_index].receipt_path)
    populate_entries(receipt_data_list[global_index])
    

# Loads the first receipt
def load_first(receipt_data_list : list[ReceiptData]):
    global global_index

    if global_index >= len(receipt_data_list):
        root.destroy()

    swap_photo(receipt_data_list[global_index].receipt_path)
    populate_entries(receipt_data_list[global_index])



# Changing the directory to be compatible with an executable
dir_path = os.path.dirname(os.path.realpath(__file__))
os.chdir(dir_path)

receipts_csv_path = 'data_intake.csv'
receipt_data_list = read_receipt_csv(receipts_csv_path)


# TKinter frame setup
csv_data = []
global_index = 0
root = ctk.CTk()
root.title('Receipt Viewer')
root.geometry('1920x1080')

main_frame = ctk.CTkFrame(root, fg_color='#333333')
main_frame.pack(anchor=tk.CENTER, fill=tk.BOTH, expand=True)

wrapper_frame = ctk.CTkFrame(main_frame, fg_color='#333333')
wrapper_frame.pack(anchor=tk.CENTER, fill=tk.Y, expand=True)


left_frame = ctk.CTkFrame(wrapper_frame, fg_color='#333333')
left_frame.pack(side=tk.LEFT)

right_frame = ctk.CTkFrame(wrapper_frame, fg_color='#333333', width=600, height=1080)

right_frame.pack(side=tk.LEFT)
right_frame.pack_propagate(False)


# Adds TKinter entry fields
supplier_label, supplier_entry = add_entry("Supplier: ")
account_label, account_entry = add_entry("Account: ")
subtotal_label, subtotal_entry = add_entry("Subtotal: ")
tax_label, tax_entry = add_entry("Tax: ")
tips_label, tips_entry = add_entry("Tips: ")
date_label, date_entry = add_entry("Date, YYYY/MM/DD: ")

# Start image
image_label = ctk.CTkLabel(right_frame, text="")
image_label.pack(expand=True)
load_first(receipt_data_list)

# Next button
next_button = ctk.CTkButton(left_frame, text="Save", command=lambda: load_next(receipt_data_list))
next_button.pack(padx=200, pady=20)

root.mainloop()