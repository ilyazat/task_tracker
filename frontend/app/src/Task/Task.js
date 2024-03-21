import React, { useEffect, useState } from 'react';
import successImage from './success.png';
import failImage from './fail.jpg';


function getCookie(cookieName) {
    const name = cookieName + "=";
    const decodedCookie = decodeURIComponent(document.cookie);
    const cookieArray = decodedCookie.split(';');

    for(let i = 0; i < cookieArray.length; i++) {
        let cookie = cookieArray[i];
        while (cookie.charAt(0) === ' ') {
            cookie = cookie.substring(1);
        }
        if (cookie.indexOf(name) === 0) {
            return cookie.substring(name.length, cookie.length);
        }
    }
    return null;
}


function TrackerPage() {
    const [authorized, setAuthorized] = useState(false);

    useEffect(() => {
        const token = getCookie("token")
        if (token) {
            setAuthorized(true);
        }
    }, []);

    return (
        <div style={{ textAlign: 'center' }}>
            {authorized ? (
                <>
                    <h1>Successfully Authorized</h1>
                    <img src={successImage} alt="Success" />
                </>
            ) : (
                <>
                    <h1>You need to authorize</h1>
                    <img src={failImage} alt="Fail" />
                </>
            )}
        </div>
    );
}

export default TrackerPage;
