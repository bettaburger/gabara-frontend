require('dotenv').config();

const clientId = process.env.SPOTIFY_CLIENT_ID;
const redirectUri = "http://127.0.0.1:8000/callback";

async function getToken(code) {
  const codeVerifier = localStorage.getItem("code_verifier");

  const body = new URLSearchParams({
    client_id: clientId,
    grant_type: "authorization_code",
    code,
    redirect_uri: redirectUri,
    code_verifier: codeVerifier,
  });

  const response = await fetch("https://accounts.spotify.com/api/token", {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body,
  });

  const data = await response.json();
  console.log("TOKEN:", data);

  localStorage.setItem("access_token", data.access_token);
}

const params = new URLSearchParams(window.location.search);
const code = params.get("code");

if (code) {
  getToken(code);
}