# Print forms

A print form puts results onto a pre-printed sheet, e.g. the lane sheets the army supplies for the obligatory programme: the shots are entered at the lane or afterwards, and the sheet is printed instead of filled in by hand. Filling in sheets by hand stays possible; print forms are optional.

## Create a form {#create}

**Print forms** (admins, in the club part of the sidebar) lists the forms of the club. **New print form** asks for a name and a paper size, then opens the editor. It has three columns:

- **Fields** on the left: everything a form can print, in groups, with an example each and a search. A click adds the field below the lowest element; dragging it onto the sheet puts it where you drop it.
  - *Own text* prints fixed text, *QR code* and *Barcode* print their value as a code, *Line* draws a line.
  - *Shooter*, *Club and event* and *Result* print one value each, e.g. the name, the total or the points of a series.
  - *Single shots* adds a whole row of shot boxes: *Shots of series 1* is one element `{{s1.shot#}}` repeated once per shot, 8 mm apart.
- The **Sheet** in the middle shows the form true to scale over the scan, with made-up values; **Show field names** shows the placeholders instead. Rulers in millimetres run along the top and the left, the zoom goes from 50 to 200 %.
- **Properties** on the right edit the selected element exactly: its **Value** (placeholders such as `{{shooter.name}}` are filled in when printing), **x** and **y** in millimetres from the top left corner, **Width** and **Height** (a width of 0 lets text run to the right edge), **Font size**, **Bold** and **Alignment**. **Repeat** prints it several times, moved by **dx** and **dy** each time; `#` in a placeholder counts 1, 2, 3, …. Below, the **Form** sets the name, the size of the sheet, the printer offset and the **Scan of the empty form**: a PNG or JPEG of the whole sheet that appears behind the sheet so the fields can be placed onto the boxes of the form; it is never printed.

**Save** stores the form; a dot on the button shows unsaved changes, and leaving the page with them asks first. With **Save automatically** (on unless you switch it off; the browser remembers it) unsaved changes are stored every minute. **PDF preview** opens the form exactly as the PDF renders it, which can differ slightly from the sheet in the browser.

## Arrange the elements {#arrange}

- **Move** an element by dragging it, **resize** it by its corner. The pink dot on the last box of a repeated element drags the spacing.
- While moving, the element **snaps** to the edges and centres of the other elements and of the sheet, and a pink guide shows the line; elsewhere it snaps to whole millimetres. Hold Alt to place it freely.
- **Select several** with Shift or Ctrl and a click, or draw a frame around them on the empty sheet. Dragging moves them together; font size, bold and alignment then apply to all of them.
- A **right click** opens the menu: *Duplicate*, *Bring to front*, *Align* (edges, centres, distribute evenly), *Centre on the sheet* and *Delete*.

| Key | Does |
|---|---|
| Arrow keys | move the selection by 0.5 mm, with Shift by 5 mm |
| Ctrl+Z / Ctrl+Y | undo / redo |
| Ctrl+A | select everything |
| Ctrl+D | duplicate |
| Delete | delete the selection |
| Esc | clear the selection |
| Ctrl+S | save |
| Ctrl+mouse wheel | zoom |

## Calibrate the printer {#calibrate}

Printers rarely print exactly where they should. Print a sheet with **Test print**, measure how far the values are off and enter it as **Printer offset right** and **Printer offset down** (negative to move left or up). The offset applies to every element; the preview ignores it.

## Use a form {#use}

Under **Configuration** → **Programmes** of the event, open a programme and choose its **Print form** in the section **Use**; this works while shooters are registered. Every completed programme then offers **Print onto form** in the results: at the desk, on the **Enter results** page and on scorer links at the lane (see [Helpers](helpers.md)). With a print agent and a printer chosen for *Print forms* under **Ranges → Printing**, the sheet goes to that printer (e.g. its manual feed), otherwise the PDF opens (see [Printing](printing.md)).

**Export** downloads a form as a file and **Import** reads it again, e.g. to share a form with another club.
