<html>
    <head>
    <title>Register</title>
    <style type="text/css">
        
    </style>
    </head>
    <body>
        <form action="/register" method="post">
            Name :<input type="text" name="name">
            Username:<input type="text" name="username">
            Email:<input type="email" name="email">
            Password:<input type="password" name="password">
            Re-Enter Password:<input type="password" name="repassword">
            <input type="hidden" name="token" value="{{.}}">
            <input type="submit" value="Register">
        </form>
    </body>
</html>