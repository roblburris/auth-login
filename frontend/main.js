async function onSignIn(googleUser) {
  const response = fetch("http://localhost:8000/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      type: "gsuite",
      jwt: googleUser.getAuthResponse().id_token,
    }),
  }).then(response => response.text()).then(
      text => document.textContent = text
  );
}

async function onSubmitLogin() {
  return;
  //   let formData = new FormData(document.getElementById("form"));
  //   const response = await fetch("http://localhost:63342/login", {
  //     method: "POST",
  //     headers: {
  //       "Content-Type": "application/json",
  //     },
  //     body: {
  //       type: "non-gsuite",
  //       //   username: formData.get("username"),
  //       //   password: formData.get("password"),
  //     },
  //   });
  //   window.alert("abc");
}
