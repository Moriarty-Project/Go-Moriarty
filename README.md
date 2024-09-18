# Go-Moriarty

 An offshoot of how I think the infosec tool Sherlock could be truly used to it's fullest
 
 The overview of this project. Moriarty will have all the actual functions and whatnot to get a users info and history, cli just handles the user interface elements. 
 If this is used in other tools, the direct library aspects of Moriarty would just be imported.
 
 The end goal here is to have an infosec tool that can semi autonomously analyze swaths of data for elements relevant to a target person. Analyzing online account findings, and local data sets. Ideally, finding cross links of common interactions across platforms to truly find all publicly available information about the person.

# The name reasoning.

Originally I couldn't figure out what to name it, but Moriarty seems the best fit. 
First of all, the famous Sherlock rival Moriarty, but further more, the name comes from an Irish name Ó Muircheartaigh [oː ˈmˠɪɾʲɪçaɾˠt̪ˠiː], which can be translated to 'navigator'. [(Source Wiki)](https://en.wikipedia.org/wiki/Moriarty_(name))

I like the idea of this being a tool to talk with, and to help you navigate



## Main items to deal with next
- [:heavy_check_mark:] Start to split project into separate packages
\- [:heavy_multiplication_x:] break the separated packages into different repos
- [:heavy_minus_sign:] Setup github actions
\- [:heavy_check_mark:] Get build action setup
\- [:heavy_check_mark:] Get basic test action for each part setup
\- [:heavy_multiplication_x:] Get race conditions tested for each package
\- [:heavy_multiplication_x:] Get core to pass remote tests
- [:heavy_multiplication_x:] Get CLI working.
\- [:heavy_multiplication_x:] write a user guide for the CLI.
- [:heavy_multiplication_x:] Change how user results are stored, and make them easier to add unique fields to, and query for information.
- [:heavy_multiplication_x:] Add methods for adding publicly listed "friends" from different sites, and ways to cross-link those items.
- [:heavy_multiplication_x:] Web graph connections between different found target users.

<!-- the signs we're using. -->
<!-- heavy_check_mark -->
<!-- heavy_minus_sign -->
<!-- heavy_multiplication_x -->