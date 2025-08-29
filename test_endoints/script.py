# test for registering and joining event

import requests
import sys

base_url = "http://localhost:8080"
user = "user"
passw = "pass"

headers = {
  'email': user,
  'password': passw,
  "username": user
}


resp = requests.request("POST", base_url + "/register", headers=headers)

if resp.status_code != 200:
    print("There is a problem ", resp.text,resp.status_code)
    sys.exit()





# payload = {}
headers = {
  'email': user,
  'password': passw
}

login_response = requests.request("POST", base_url + "/login", headers=headers)

# print(response.text)



if login_response.status_code != 200:
    print("eror while login ", login_response.text)
    sys.exit()

body = login_response.json()
id = body.get("id")
level = body.get("level")
token = body.get("token")





# import requests

url = "http://localhost:8080/join"

payload = {}
headers = {
  'Authorization': f'Bearer {token}',
  'eventNo': '4',
  'userId': str(id),
  'level': str(level)
}

response = requests.request("POST", base_url + "/join", headers=headers)

print(response.text)
