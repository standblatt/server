# How printing works

A shooting festival prints a lot: a lane sheet and a control strip for every registration, result sheets, wreath cards, forms. Standblatt is a web app, and that shapes how printing works.

## Every document is a PDF {#pdf}

The server renders every document as a PDF, in the language of the event: lane sheets, control strips, result sheets, wreath cards, print forms, certificates and rankings. The same PDF is printed at the range, opened in the browser or shared from a phone, so a document looks the same wherever it comes out. Rankings are only downloaded, as PDF or CSV, and never sent to a printer of the range.

## Why a print agent {#agent}

A browser cannot print on its own: every print asks for a dialog, a printer and a click. At a busy desk that is too slow, and a lane sheet should come out of the printer next to the desk while the control strip comes out of the label printer.

The **print agent** solves this. It is a small program on a computer next to the printers: the Standblatt Print Agent app on Windows or a Mac, or a headless agent on a Linux computer such as a Raspberry Pi. When someone presses a Print button, the server puts the PDF into a queue; the agent fetches it and prints it on the chosen printer without asking. The agent connects to Standblatt, not the other way round, so nothing at the range has to be reachable from the internet.

Each print job is confirmed by the agent. A job that is not confirmed within two minutes is handed out again, up to three attempts; **Recent jobs** shows what happened.

## Agents and printers belong to a range {#range}

Printers stand at a range, not in an event, so the print agent and the printer settings belong to the [range](../getting-started/concepts.md#ranges-and-lanes). Every event on that range uses them without any set-up of its own. A registration's documents print on the range of its lane; an event on several ranges therefore prints at the range where the shooter shoots. One computer can serve several ranges, e.g. *Pistol 25 m* and *Pistol 50 m* in the same hall.

The settings choose a printer **per document**, because the documents need different printers: the control strip goes to a label printer, the lane sheet to an A5 tray, wreath cards and print forms to a manual feed with the pre-printed sheets. The printer options the agent reports (tray, paper size, two-sided and so on) are stored with the choice.

An agent is connected either with a token that is copied into its settings or, with the app, by a code shown on the computer and confirmed in Standblatt (see [Printing](../user-guide/printing.md#connect)). Either way the server keeps only a fingerprint of the token, never the token itself.

## When there is no printer {#fallback}

A document goes to the printer only if its range has a printer chosen for that document and an agent of the range was online within the last minute. Otherwise the PDF opens instead: in a new tab on a computer, in the share sheet on a phone or tablet (print, save to files, send), so the app stays open. Nothing is lost when the agent is off; the desk just prints by hand.

The same holds for events without a range: there is no agent to send to, so every document opens as PDF.

## Print forms on pre-printed sheets {#forms}

Some documents are not designed by Standblatt but supplied on paper: wreath cards, the lane sheets the army supplies for the obligatory programme, a club's own certificates. A **print form** describes such a sheet: where on the paper each value goes, in millimetres from the top left corner. Printing a form puts only the values onto the sheet that is already in the printer.

A form has a **data source**: the results of a programme, a registration or a rank. The source decides which values the form can print and where it is used: forms of results are printed after a programme is completed or as wreath cards, forms of a registration as control strips or participant cards, forms of a rank as certificates from the rankings. The wreath card and the control strip are themselves built-in forms that a club can copy and adjust.

A scan of the empty sheet can lie behind the form in the editor, so the fields can be placed onto its boxes. The scan is for placing only and is never printed.

## Calibration in principle {#calibration}

Printers rarely print exactly where they should: paper feeds and trays shift the print by a millimetre or two. A form that is right on screen can therefore land beside the boxes on paper. The **printer offset** of the form moves every element by the same amount to make up for it: print a test sheet, measure how far the values are off and enter the difference. The offset is part of the form, so a form printed on another printer that shifts differently needs a copy with its own offset. The preview ignores the offset and shows the form as designed. The steps are in [Print results onto the army's lane sheets](../how-to/army-sheet.md#calibrate).
