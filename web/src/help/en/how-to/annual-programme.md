# Run an annual programme over the season

You need this when members shoot a programme as often as they like over the season and a club championship ranks the best three results of each shooter. It needs the permission *Configure events* (see [Roles](../user-guide/users.md#roles)); entering the results needs *Run events*.

## Create the event {#event}

1. Under **Events**, click **New event** and choose **Annual programme or training**. It switches on the modules *Attempts* and *Classifications*; choose **Club championship** instead if the ranking should also run on a kiosk display.
2. Enter the **Name**, e.g. *Annual programme 2027*, and as **Date** the first day of the season. Its year sets the age categories.
3. Tick the range under **Ranges**, or none if the results come on paper only (see [Events without lanes](../user-guide/events.md#no-lanes)).
4. Switch on **Period (season)**, enter the **First day** and the **Last day**, and under **Shooting times** tick the weekdays of the club's training with their times. No weekday ticked means every day of the period.
5. Click **Create**.

## Add the programme {#programme}

1. Open **Configuration** → **Programmes** and click **New programme** (or **Add from …** to take it from the library, a file or an earlier event).
2. Enter the name and the series.
3. Under **Attempts**, tick **Unlimited** and set **Counting attempt** to **The best**.
4. Under **Prices (CHF)**, tick every age category that may shoot and enter the price. The second field is the price of every further attempt: empty means the same price, 0 means the programme is paid once.
5. Click **Create**.

## Rank the best three attempts {#classification}

1. Open **Configuration** → **Classifications** and click **New classification**.
2. Enter a **Name**, e.g. *Club championship*, choose the annual programme as **Programme 1** and set **Counting attempts** to 3.
3. Switch on **Rank only shooters with all counting results** if only shooters with three results should get a rank; the others are then listed below the ranking.
4. Add **Classes** if juniors and adults are ranked apart, e.g. *Juniors* with U17 to U21 and *Elite* with the others.
5. Click **Create**.

All settings of the editor are described in [Classifications](../user-guide/classifications.md).

## Enter the results during the season {#results}

1. Open the event and click **Enter results** in the sidebar.
2. Choose the **Programme**, type a licence number or a name into **Shooter** and press ++enter++.
3. Type the series and press ++enter++. Every result of a shooter is a new attempt.

Results from the scoring system can also be imported as a file (see [Import results from the scoring system](import-results.md)).

**Rankings** shows the programme with the best attempt of each shooter (*2/5*: the second of five attempts counts) and the club championship with the sum of the best three.

## End the season {#close}

When the last results are in and checked, open **Configuration** → **Copy & close** and click **Close event**. What closing changes is explained in [Closing an event](../explanation/closing.md); next season starts with [Take over last year's event](next-year.md).
