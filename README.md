# VT-Scanner
File &amp; URL scanning Using VirusTotal APIs

__NOTE:__ **This Scanner Supports only File Scanning. URL Scanning will be available soon** 

__USAGE:__
>./vtScanner -f \<filename\> -a \<apikey\>

>go run vtScanner.go -f \<filename\> (IF API KEY IS INTEGRATED WITH CONFIGURATION)

__EXAMPLE:__
>./vtScanner -f /home/test.txt -a xxxxxxxxxxxxxxxxxxx

>go run vtScanner.go -f /home/test.txt (IF API KEY IS INTEGRATED WITH CONFIGURATION)


__HOW TO CREATE API KEY__
>1. Go to https://www.virustotal.com/gui/
>2. Click on *Sign up* button present at top-right corner of page
>3. Fill all the required information & click on *Join us* button
>4. Come back to https://www.virustotal.com/gui/ and click on *Sign in* button which is at left side of *Sign up* button
>5. After successfull login, Click on your avatar at top-right corner & then click on *API key*
>6. Copy You API KEY(long hash) and use save


__HOW TO INTEGRATE API KEY IN CODE(NOT RECOMMENDED)__
>1. Open VT-Scanner/scanmod/config/config.go
>2. At line number 9, Replace *xxxxxxxxxxxxxxxxxxxxxx* with your api key

>__Note:- To run scanner after editing the config file, you will need golang installed on your machine__

**Thanks**
>Avinash Ghadshi

