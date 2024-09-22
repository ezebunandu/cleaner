# Cleaner

When you take a screenshot on MacOS, it default saves to your desktop with a naming convention like `Screenshot 2023-12-13 at 10.46.37 PM.png` Cleaner is a command line tool that takes a source directory and a target and moves all screenshot files from the source to the target.

The screenshots are saved in subfolders within the target using the `date` from the filename. This can be run as a cron job to continously declutter your desktop but also make it easy for you to review screenshots from a particular date if needed.

```usage: cleaner <source> <target>```
