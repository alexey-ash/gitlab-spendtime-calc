## Small cli app for check tasks spend time 

To run, it is necessary in the main.go fill in constants:
* __gitlabURL__ - Gitlab API URL (https://gitlab.com/api/v4)
* __gitlabToken__ - AccessToken from Project where tasks 
* __projectID__ - ID of project with tasks
* __userName__ - Your username

## Result
```
[~]: gitlab-time
+----------------------------------------------------------------+-------+
|                          TASK                                  | HOURS |
+----------------------------------------------------------------+-------+
| AP-74 Implement authentication module for web application      |     2 |
| AP-66 Create REST API for mobile app using Node.js and MongoDB |     5 |
| AP-29 Optimize database queries for faster performance         |     2 |
+----------------------------------------------------------------+-------+
Total spend time today: 9h
```



