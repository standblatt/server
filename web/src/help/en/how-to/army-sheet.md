# Print results onto the army's lane sheets

You need this when the shooters' results must end up on pre-printed sheets, e.g. the lane sheets the army supplies for the obligatory programme, and you would rather print them than fill them in by hand. It needs the permissions *Print forms*, *Configure events* and *Ranges and printers* (see [Roles](../user-guide/users.md#roles)).

Before you start, scan an empty sheet as PNG or JPEG and keep a few spare sheets for test prints.

## 1. Prepare the programme {#programme}

1. Check that **Club modules** in the club part of the sidebar has *Print forms* switched on (it is on unless someone switched it off).
2. Under **Configuration** → **Programmes** of the event, open the programme and tick **Single shots** for every series whose shots the sheet has boxes for. Do this before the first registration: once shooters are registered, the series are fixed.
3. **Save** the programme.

## 2. Create the form {#form}

1. Open **Print forms** in the club part of the sidebar and click **New print form**.
2. Enter a name (e.g. *Obligatory programme 300 m*), choose **Prints** *Results of a programme* and the paper size closest to the sheet, then **Create**. The editor opens.
3. In **Form** on the right, set **Width (mm)** and **Height (mm)** to the size of the sheet, then **Upload a scan** under **Scan of the empty form**. The scan shows behind the sheet while editing; it is never printed.
4. From **Fields** on the left, drag the shooter's fields (e.g. **Last name**, **First name**, **Year of birth**, **License number**) onto their boxes on the scan.
5. From the group **Single shots**, drag **Shots of series 1** onto the first shot box of that series. It is a row of boxes, 8 mm apart; drag the pink dot on its last box until every value sits in its box, or set **dx (mm)** under **Properties**. Repeat for the other series.
6. From **Result**, add the series totals (**Series 1** …) and the **Total** where the sheet has boxes for them.
7. **Save**.

![Print forms](../images/screenshots/print-forms.png){ .screenshot }

How to move, align and snap elements is in [Print forms](../user-guide/print-forms.md#arrange), every field in [Form placeholders](../reference/form-placeholders.md).

## 3. Calibrate the printer {#calibrate}

1. Put an empty sheet into the printer at the range.
2. Click **Test print**: the saved form opens as PDF with made-up values and without the scan. Print it at actual size (no fit to page).
3. Measure how far the values are off their boxes and enter it as **Printer offset right (mm)** and **Printer offset down (mm)** under **Form** (negative to move left or up).
4. **Save** and repeat from step 1 until the values sit in the boxes.

## 4. Assign the form and the printer {#assign}

1. Under **Configuration** → **Programmes** of the event, open the programme and choose the form in the section **Print form**. It is saved at once; this works while shooters are registered.
2. Under **Ranges → Set up the printers** of the range, choose the printer for **Print forms (pre-printed)**, e.g. with its manual feed as tray, and **Save** (see [Connect a printer](connect-printer.md)).

## 5. Print the sheets {#print}

Once all series of the programme are in, its results show **Print onto form**: in the results at the desk and on scorer links at the lane (see [Helpers](../user-guide/helpers.md)). On **Enter results**, every line under *Entered now* has its own print button. Put a sheet into the printer and click it. Without an online print agent, the PDF opens and is printed from the browser.
