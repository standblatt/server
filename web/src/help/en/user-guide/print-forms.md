# Print forms

A print form puts results onto a pre-printed sheet, e.g. the lane sheets the army supplies for the obligatory programme: the shots are entered at the lane or afterwards, and the sheet is printed instead of filled in by hand. Filling in sheets by hand stays possible; print forms are optional and can be switched off under **Club modules** (see [Club modules](events.md#club-modules)).

## Create a form {#create}

**Print forms** (admins, in the club part of the sidebar) lists the forms of the club. **New print form** asks for a name, what the form **Prints** (see [Data source](#source)) and a paper size, then opens the editor. It has three columns:

- **Fields** on the left: everything a form can print, in groups, with an example each and a search. A click adds the field below the lowest element; dragging it onto the sheet puts it where you drop it.
  - *Own text* prints fixed text, *QR code* and *Barcode* print their value as a code, *Line* draws a line.
  - *Shooter*, *Club and event* and *Result* print one value each, e.g. the name, the total or the points of a series.
  - *Single shots* adds a whole row of shot boxes: *Shots of series 1* is one element `{{s1.shot#}}` repeated once per shot, 8 mm apart.
- The **Sheet** in the middle shows the form true to scale over the scan, with made-up values; **Show field names** shows the placeholders instead. Rulers in millimetres run along the top and the left, the zoom goes from 50 to 200 %.
- **Properties** on the right edit the selected element exactly: its **Value** (placeholders such as `{{shooter.name}}` are filled in when printing), **x** and **y** in millimetres from the top left corner, **Width** and **Height** (a width of 0 lets text run to the right edge), **Font size**, **Bold**, **Italic**, **Wrap text** (the text breaks into lines within the width, each as high as the height) and **Alignment**. **Repeat** prints it several times, moved by **dx** and **dy** each time; `#` in a placeholder counts 1, 2, 3, …. Below, the **Form** sets the name, the size of the sheet, the printer offset and the **Scan of the empty form**: a PNG or JPEG of the whole sheet that appears behind the sheet so the fields can be placed onto the boxes of the form; it is never printed.

**Save** stores the form; a dot on the button shows unsaved changes, and leaving the page with them asks first. With **Save automatically** (on unless you switch it off; the browser remembers it) unsaved changes are stored every minute. **PDF preview** opens the form exactly as the PDF renders it, which can differ slightly from the sheet in the browser.

## Data source {#source}

What a form **Prints** decides its fields and where it is printed:

| Prints | Fields | Printed from |
|---|---|---|
| Results of a programme | shooter, club, event, series, single shots, total, award and its value | the results of a programme (see [Use a form](#use)) or as the event's wreath card |
| A registration | shooter, club, event, lane, programmes, passes, QR code of the registration | as the event's control strip |
| A rank | shooter, club, event, rank, total, ranking, class | certificates from the rankings |

The source can be changed in the editor under **Form**; fields the new source does not know print nothing.

## Built-in forms {#builtin}

Standblatt prints the wreath card and the control strip with **built-in forms**, listed below the forms of the club. **Copy** turns one into a form of the club that can be changed like any other, e.g. to move the fields onto another pre-printed wreath card or to add the club's logo. An event then chooses it under **Configuration** → **Modules**: *Print form of the wreath cards* in the module *Awards*, *Print form of the control strips* in the module *Control strip*. Without a choice, or when the chosen form is deleted, the built-in form prints.

![Print forms with the built-in forms and the image library](../images/screenshots/print-forms.png){ .screenshot }

## Registration documents and certificates {#certificates}

A form that prints **A registration** besides the event's control strip, e.g. a participant card or a receipt, appears in the registration with a button of its own name; it prints on the printer of *Print forms*.

A form that prints **A rank** makes certificates: **Certificates** in the header of every ranking and classification on **Rankings** asks for the form, the class of a classification and the ranks (1 to 3 unless changed), then prints one page per ranked shooter, or opens them as one PDF. Certificates go to the printer of *Certificates* under **Ranges → Printing**. The fields of a rank are the rank, the total, the ranking (the programme or classification) and the class; a programme's certificate knows its results too.

![Certificates from a ranking](../images/screenshots/certificates.png){ .screenshot }

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

## Images {#images}

Logos, letterheads and other pictures come from the club's image library below the list of **Print forms**: **Upload image** adds a PNG or JPEG of at most 2 MB, named after its file. In the editor, the group **Images** of the fields lists them; a click or a drag puts one onto the sheet, 30 mm wide with the image's own proportions. Its properties choose another **Image**, its position and size, and **Keep the aspect ratio** sets the height to match the width again. An image can be repeated like any other element.

An image can be deleted only while no form shows it. **Export** puts the images of a form into its file, and **Import** adds them to the library, reusing an image of the same name the club already has.

## Calibrate the printer {#calibrate}

Printers rarely print exactly where they should. Print a sheet with **Test print**, measure how far the values are off and enter it as **Printer offset right** and **Printer offset down** (negative to move left or up). The offset applies to every element; the preview ignores it.

## Use a form {#use}

Under **Configuration** → **Programmes** of the event, open a programme and choose its **Print form** in the section of the same name; this works while shooters are registered. Every completed programme then offers **Print onto form** in the results: at the desk, on the **Enter results** page and on scorer links at the lane (see [Helpers](helpers.md)). With a print agent and a printer chosen for *Print forms* under **Ranges → Printing**, the sheet goes to that printer (e.g. its manual feed), otherwise the PDF opens (see [Printing](printing.md)).

**Export** downloads a form as a file, with its images, and **Import** reads it again, e.g. to share a form with another club.
