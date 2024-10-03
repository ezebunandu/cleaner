# Cleaner

`cleaner` takes a source directory and a target and moves all screenshot files from the source to the target.

To install and use cleaner, you need to have Go >= 1.23.1 installed. Simply run the following command:

`go install github.com/ezebunandu/cleaner/cmd/cleaner@v0.1.2`

When you take a screenshot on MacOS, it default saves to your desktop with a naming convention like `Screenshot 2023-12-13 at 10.46.37 PM.png`. When you run cleaner, it would declutter your desktop by moving all screenshot files to the target directory and organize them into subfolders by date, making it easy for you to go back and find a particular screenshot, should you need to!

I have it running as a cron job!.

```usage: cleaner <source directory> <target directory>```
