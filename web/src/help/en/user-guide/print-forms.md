# Print forms

A print form puts results onto a pre-printed sheet, e.g. the lane sheets the army supplies for the obligatory programme: the shots are entered at the lane or afterwards, and the sheet is printed instead of filled in by hand. Filling in sheets by hand stays possible; print forms are optional.

## Create a form {#create}

**Print forms** (admins, in the club part of the sidebar) lists the forms of the club. **New print form** asks for a name and a paper size, then opens the editor:

- **Width** and **Height** of the sheet in millimetres.
- **Scan of the empty form**: a PNG or JPEG of the whole sheet. It appears behind the **Preview** so the fields can be placed onto the boxes of the form; it is never printed.
- **Elements**: what is printed where. **x** and **y** are millimetres from the top left corner, **Width** and **Height** the size (a width of 0 lets text run to the right edge), **Size** the font size in points.
  - *Text* prints its **Value**, in which placeholders are replaced, e.g. `{{shooter.name}}` or `{{total}}`.
  - *QR* and *Barcode* print their value as a code, *Line* draws a line from x/y by width/height.
  - **Repeat** prints the element several times, moved by **dx** and **dy** each time; `#` in a placeholder counts 1, 2, 3, …. One element `{{s1.shot#}}` with repeat 10 and dx 8 fills the ten shot boxes of series 1.
- **Placeholders** lists every value with an example; a click adds it to the element edited last.

**Save** refreshes the preview, which shows made-up values.

## Calibrate the printer {#calibrate}

Printers rarely print exactly where they should. Print a sheet with **Test print**, measure how far the values are off and enter it as **Printer offset right** and **Printer offset down** (negative to move left or up). The offset applies to every element; the preview ignores it.

## Use a form {#use}

On the *Programmes* tab of the event's **Setup**, choose the **Print form** of a programme; this works while shooters are registered. Every completed programme then offers **Print onto form** in the results: at the desk, on the **Enter results** page and on scorer links at the lane (see [Helpers](helpers.md)). With a print agent and a printer chosen for *Print forms* under **Ranges → Printing**, the sheet goes to that printer (e.g. its manual feed), otherwise the PDF opens (see [Printing](printing.md)).

**Export** downloads a form as a file and **Import** reads it again, e.g. to share a form with another club.
