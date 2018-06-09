<html>
    <head>
        <title>Home Page</title>
        <style type="text/css">
             ul {
                    list-style-type: none;
                    margin: 0;
                    padding: 0;
                    overflow: hidden;
                    background-color: #333;
                }

                li {
                    float: left;
                    border-right:1px solid #bbb;
                }

                li:last-child {
                    border-right: none;
                }

                li a {
                    display: block;
                    color: white;
                    text-align: center;
                    padding: 14px 16px;
                    text-decoration: none;
                }

                li a:hover:not(.active) {
                    background-color: #111;
                }

                .active {
                    background-color: #00B7D5;
                    color: #020101;
                }
        </style>
    </head>
    <body>
        <form>
            <ul>
                <li> <a class="active" href="/login">Login</a> </li>
                <li> <a href="/register">Register</a> </li>
                <li style="float:right" > <a href="#">About</a> </li>
            </ul> 
        </form>
    </body>
</html>