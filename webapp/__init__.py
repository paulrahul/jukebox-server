import os

from flask import Flask, g, render_template, request, url_for, session
from werkzeug.utils import redirect

import radio5
import spotify

def create_app():
    app = Flask(__name__)
    app.config.from_mapping(
        SECRET_KEY="rahulp0aul",
        CLIENT_ID = os.environ.get("CLIENT_ID"),
        CLIENT_SECRET = os.environ.get("CLIENT_SECRET"),
        UNAME = os.environ.get("UNAME"),
        PASSWD = os.environ.get("PASSWD"),
        CONTRIBUTOR_ID = os.environ.get("CONTRIBUTOR_ID")
    )

    @app.route("/")
    def index():
        return render_template("index.html")

    @app.route("/consolidate")
    def consolidate():
        # radio5_likes = get_radio5().fetch_radio5_likes(
        #     app.config['CONTRIBUTOR_ID'])

        spotify_user = get_spotify().get_current_user(session['spotify_at'])

        return render_template(
            "list.html",
            spotify_user=spotify_user)
            # radio5_likes=get_radio5().import_likes(radio5_likes))

    @app.route("/login")
    def login():
        return redirect(get_spotify().auth_url())

    @app.route("/auth_callback")
    def auth_callback():
        try:
            code = request.args.get("code")
        except KeyError:
            raise ValueError("Error in authorization. No code found")
            return redirect(url_for('index'))            
            
        ret, at, rt = get_spotify().get_auth_access_token(code)
        if ret == False:
            raise ValueError("Error in authorization. No code found")
            return redirect(url_for('index'))

        session['spotify_at'] = at
        session['spotify_rt'] = rt
            
        return redirect(url_for('consolidate'))

    def get_spotify():
        if 'spotify' not in g:
            g.spotify = spotify.Spotify(app.config["CLIENT_ID"],
                                        app.config["CLIENT_SECRET"])

        return g.spotify

    def get_radio5():
        if 'radio5' not in g:
            g.radio5 = radio5.Radio5Session()
            g.radio5.login_radio5(app.config['UNAME'], app.config['PASSWD'])

        return g.radio5

    return app

if __name__ == "__main__":
    create_app().run(port=8080)