# Swaggy-New-Youtube-Client
youtube client because youtube is fucking stupid.

step 0: Have ffmpeg installed
I always think this is installed by default if it ain't ```sudo apt install ffmpeg``` I hope whatever version your distros repository gives you works. If you're on windows. Use WSL debian idk

Step 1: download the latest version of yt-dlp
Guide here:
https://github.com/yt-dlp/yt-dlp#installation

Step 2: move that yt-dlp python file to $path and make it an executable so you can run it by typing "yt-dlp" in your command line
I downloaded it manually and from my downloads folder ran: 
```mv yt-dlp /usr/local/bin/yt-dlp```

Then from /usr/local/bin we run:
```sudo chmod +x /usr/local/yt-dlp```

Totally unnecassary to be in the same dir as the file if we specify the filepath anyway. Whatever gets the job down. Now from your home folder (or any dir I think) you should be able to run and update yt-dlp!

step 3: Install MPV 
guide here if you need it: https://github.com/mpv-player/mpv

I think MPV is installed by default on your distro. It should atleast be in your package manager so if you're on debian or ubuntu ```sudo apt install mpv``` should work.
I really can't rememeber if I even had to install mpv or not. This isn't a very good guide.

step 4: Download and run our awful youtube client. If you want to search for and watch a video. It works. Not very well. But it works.

