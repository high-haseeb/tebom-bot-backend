async function login() {
    const response = await fetch("https://portal.acente365.com/account/login", {
        headers: {
            "accept": "*/*",
            "accept-language": "en-US,en;q=0.9",
            "cache-control": "no-cache",
            "content-type": "application/x-www-form-urlencoded; charset=UTF-8",
            "pragma": "no-cache",
            "priority": "u=1, i",
            "sec-ch-ua": "\"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"",
            "sec-ch-ua-mobile": "?0",
            "sec-ch-ua-platform": "\"Linux\"",
            "sec-fetch-dest": "empty",
            "sec-fetch-mode": "cors",
            "sec-fetch-site": "same-origin",
            "x-requested-with": "XMLHttpRequest",
            "Referer": "https://portal.acente365.com/account",
            "Referrer-Policy": "strict-origin-when-cross-origin"
        },
        body: "loginModel%5BUsername%5D=rizacenkercivelek&loginModel%5BPassword%5D=9R9pKrkI5dRhgJEyND0LSg%3D%3D&loginModel%5BRedirectUrl%5D=&loginModel%5BMfaControl%5D=true&loginModel%5BMfaCode%5D=830808&loginModel%5BGuid%5D=186b5417-793b-4358-85b1-a9ad5b7d5471",
        method: "POST",
        credentials: "include" // Ensures cookies are included in the request
    });

    // Log the full response headers
    console.log("Response Headers:");
    for (let [key, value] of response.headers.entries()) {
        console.log(`${key}: ${value}`);
    }

    // Log the Set-Cookie header (may not be accessible due to CORS restrictions)
    console.log("Set-Cookie Header:", response.headers.get("set-cookie"));

    // Log response body (if needed)
    const responseBody = await response.text();
    console.log("Response Body:", responseBody);
}

login();

