const clientId = "105507cca7d448d7a4f7ce2108cc5f37";
const redirectUri = "'https://bettaburger.github.io/gabara/docs/callback.html'";

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