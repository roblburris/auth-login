async function onSignIn(googleUser) {
  let profile = googleUser.getBasicProfile();
  console.log("ID: " + profile.getId()); // Do not send to your backend! Use an ID token instead.
  console.log("Name: " + profile.getName());
  console.log("Image URL: " + profile.getImageUrl());
  console.log("Email: " + profile.getEmail()); // This is null if the 'email' scope is not present.
  console.log(googleUser.getAuthResponse().id_token);

  const response = fetch("http://localhost:63342/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      type: "gsuite",
      jwt: googleUser.getAuthResponse().id_token,
    }),
  });
  window.alert("a");
  console.log(googleUser.getAuthResponse().id_token);
  window.location = "signed-in.html";
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
