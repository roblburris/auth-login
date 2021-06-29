async function onSignIn(googleUser) {
    let profile = googleUser.getBasicProfile();
    console.log('ID: ' + profile.getId()); // Do not send to your backend! Use an ID token instead.
    console.log('Name: ' + profile.getName());
    console.log('Image URL: ' + profile.getImageUrl());
    console.log('Email: ' + profile.getEmail()); // This is null if the 'email' scope is not present.
    console.log(googleUser.getAuthResponse().id_token);

    const response = await fetch('http://localhost:63342/login', {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
        body: {
            "type": "gsuite",
            "jwt": googleUser.getAuthResponse().id_token
        }
    });
    console.log(response.json())
    window.location = "signed-in.html";
}

function onSubmitLogin() {
    console.log("yolo");
    let formData = new FormData(document.getElementById('form'))
    window.alert(formData.toString());
}

