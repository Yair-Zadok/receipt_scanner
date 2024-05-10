This is a prototype of a QuickBooks receipt scanner which pre-populates fields using AI then gives a GUI for a human to check over the AI. The scope of this project was to test the usefulness of the concept in a quick fashion.

A sample of how the GUI looks can be found at 'HowItLooks.png'.


The program is architected such that:

get_data.py produces a file called data_intake.csv which then allows display.py to display a GUI
with receipt fields pre-populated with data from data_intake.csv.

This was done in two components to allow for the AI prompt retrievals to happen asynchronously from the GUI validation step.

When the user completes validating the AI using the GUI a file called output.csv is made which can be imported into QuickBooks through their import data feature.








Copyright 2024, Yair Zadok, All rights reserved.

