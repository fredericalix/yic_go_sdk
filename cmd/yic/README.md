# yic cli managing yourITcity tokens

## Install

    go get -u github.com/youritcity/go-sdk/yic

## Usage

### See every command

    yic -h

### Login or signup

    yic login your@email.com

    yic signup your@email.com

To use other command of the cli you must set the env variable `YOURITCITY_APPTOKEN` with the token given by the signup or login command. Don't forget to validate the token sent you by email.

## Roles

List the possibles authorization roles

    yic roles

## List token

    yic token list

## Create a new token

Create a new token with a given authorization role. It works exaclty like a login but it is already validated without the sending an email.

    yic token create sensor

## Revoke a token

    yic token revoke 34LXbabCQSGnKJqxk5oi9w
