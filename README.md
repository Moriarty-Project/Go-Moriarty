# Go-Moriarty

 an offshoot of how I think the infosec tool Sherlock could be truly used to it's fullest
 
 the overview of this project. Sherlock will have all the actual functions and whatnot to get a users info and history, cli just handles the user interface elements. 
 If this is used in other tools, sherlock would just be imported.
 I will have to come up with a different name than sherlock... at the very least, to have it differentiated from the python application that inspired it...

# The name reasoning.

Originally I couldn't figure out what to name it, but Moriarty seems the best fit. 
First of all, the famous Sherlock rival Moriarty, but further more, the name comes from an Irish name Ó Muircheartaigh [oː ˈmˠɪɾʲɪçaɾˠt̪ˠiː], which can be translated to 'navigator'. [(Source Wiki)](https://en.wikipedia.org/wiki/Moriarty_(name))

I like the idea of this being a tool to talk with, and to help you navigate



## Main items to deal with next
- improve performance. 
    - the run time process takes far too long, with too many go routines.
- setup cli
    - I want the CLI to solve delay with usability. I want you to be able to set it, and forget it. Moriarty should be able to semi automatically start processing out information from a user.
- improve usability
    - beyond just what I want improved in speed, or with the CLI, I want better user controls for what you're running.