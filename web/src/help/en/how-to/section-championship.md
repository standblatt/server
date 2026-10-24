# Rank a section championship over several programmes

You need this when the ranking adds up several programmes of one event, e.g. a section championship as the sum of *Obligatory*, *Field match* and *Section stich*. All programmes must belong to the same event. It needs the permission *Configure events* (see [Roles](../user-guide/users.md#roles)).

1. Open the event and go to **Configuration** → **Modules**. Switch on **Classifications** unless the event's kind already did.
2. Under **Programmes**, make sure every programme of the championship exists (see [Programmes](../user-guide/programmes.md)).
3. Open **Classifications** in the configuration and click **New classification**.
4. Enter the **Name**, e.g. *Section championship*.
5. Choose the first programme as **Programme 1**, then **Add programme** for each further one. Leave **Counting attempts** at 1 for every programme.
6. Switch on **Rank only shooters with all counting results**. Shooters who are still missing a programme are listed below the ranking without a rank.
7. Leave **Count only the best programmes** at 0, so that every programme counts. For *the best two of three*, enter 2.
8. Check the **Tie-break**: the default is the SSV order; the arrows change the order of the criteria.
9. Add **Classes** if the championship ranks age groups apart, e.g. *Juniors* with U17 to U21 and *Elite* with the others. Without classes, everybody is in one ranking.
10. Switch on **Show on kiosk displays** if the ranking should run in the hall (only with the module *Kiosk displays*).
11. Click **Create**.

The results are entered as usual for each programme. **Rankings** shows the championship after the programme rankings, one table per class; its **PDF** and **CSV** are the report *Classifications*, and **Certificates** prints certificates for the best ranks onto a [print form](../user-guide/print-forms.md#certificates).

How totals and ties are decided is explained in [Rankings](../explanation/rankings.md).
