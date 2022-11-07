import base64
import json
import requests
import urllib

class Spotify:
    def __init__(self, client_id, client_secret):
        self.BASE_URL = "https://api.spotify.com/v1"
        self.AUTH_BASE_URL = "https://accounts.spotify.com"
        self.REDIRECT_URI = "http://localhost:8080/auth_callback"

        self._client_id = client_id
        self._client_secret = client_secret

        self._access_token = None
        self._refresh_token = None

        # self._user = None

    def auth_url(self):
        params = {
            'client_id': self._client_id,
            'response_type': 'code',
            'redirect_uri': self.REDIRECT_URI
        }

        return self.AUTH_BASE_URL + "/authorize?" + urllib.parse.urlencode(params)

    def _get_encoded_client_creds(self):
        # encoding = 'utf-8'
        client_cred_str = self._client_id + ':' + self._client_secret
        return base64.b64encode(client_cred_str.encode()).decode()

    def get_auth_access_token(self, code):
        headers = {
            'Authorization': 'Basic ' + 
                self._get_encoded_client_creds(),
            'Content-Type' : 'application/x-www-form-urlencoded'
        }

        res = requests.post(
            self.AUTH_BASE_URL + '/api/token', 
            data = {
                'grant_type': 'authorization_code',
                'code': code,
                'redirect_uri': self.REDIRECT_URI
            },
            headers = headers
        )

        if res.status_code != 200:
            print("Could not fetch Spotify access token due to: " + res.reason)
            return False, None, None

        self._access_token = json.loads(res.text)['access_token']
        self._refresh_token = json.loads(res.text)['refresh_token']
        return True, self._access_token, self._refresh_token

    def get_refresh_token(self):
        headers = {
            'Authorization': 'Basic ' + 
                self._get_encoded_client_creds(),
            'Content-Type' : 'application/x-www-form-urlencoded'
        }

        res = requests.post(
            self.AUTH_BASE_URL + '/api/token', 
            data = {
                'grant_type': 'refresh_token',
                'refresh_token': self._refresh_token
            },
            headers = headers
        )

        if res.status_code != 200:
            print("Could not fetch Spotify refresh token due to: " + res.reason)
            return False

        self._refresh_token = json.loads(res.text)['access_token']
        return True

    def get_current_user(self, at):
        headers = {
            'Authorization': 'Bearer ' + at,
            'Content-Type' : 'application/json'
        }

        res = requests.get(
            self.BASE_URL + '/me', 
            headers = headers)

        if res.status_code != 200:
            print("Could not fetch Spotify user details: " + str(res.status_code) + ", " + res.text)
            return False

        user_data = json.loads(res.text)
        return user_data['display_name']

    def _get_client_credentials_access_token(self):
        headers = {
            'Authorization': 'Basic ' + 
                self._get_encoded_client_creds(),
            'Content-Type' : 'application/x-www-form-urlencoded'
        }

        res = requests.post(
            self.AUTH_BASE_URL + '/api/token', 
            data = {'grant_type': 'client_credentials'},
            headers = headers)

        if res.status_code != 200:
            print("Could not fetch Spotify access token due to: " + res.reason)
            return False

        self._access_token = json.loads(res.text)['access_token']
        return True

if __name__ == "__main__":
    file_name = "./props.txt"

    props = {}
    with open(file_name) as f:
        for l in f.readlines():
            entry = l.split('=')
            props[entry[0]] = entry[1].strip('\n')

    spotify = Spotify(
        props['spotify_client_id'], props['spotify_client_secret'])
    

