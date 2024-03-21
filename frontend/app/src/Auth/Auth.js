import React, {useEffect, useState} from 'react';
import './Auth.css';

function Notification({ message, onClose }) {
    // Automatically close the notification after 3 seconds
    useEffect(() => {
        const timer = setTimeout(() => {
            onClose();
        }, 3000);

        return () => clearTimeout(timer);
    }, [onClose]);

    return (
        <div className="notification">
            <p>{message}</p>
        </div>
    );
}

function eraseCookie(name) {
    document.cookie = name + '=; Max-Age=0'
}


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

function AuthPage() {
    const [username, setUsername] = useState('');
    const [role, setRole] = useState('user');
    const [showNotification, setShowNotification] = useState(false);

    const handleSignup = async () => {


        try {
            const userData = {
                name: username,
                role: role
            };

            const response = await fetch('http://localhost:5554/signup', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(userData)
            });
            const body = await response.json()
            console.log(response);

            if (!response.ok) {
                throw new Error('Signup failed');
            }

            setShowNotification('Signup successful123!');
            setShowNotification(true);
            // Reset input fields
            setUsername('');
            setRole('admin');

            document.cookie = `token=${body["access_token"]};expires=${body["expires_in"]}; path=/`;

            localStorage.setItem("token", body["access_token"])
            console.log(body)

        } catch (error) {
            console.error('Error:', error);
            setShowNotification('Signup failed. Please try again.');
            setShowNotification(true);
        }
    };

    const handleLogin = async () => {
        try {
            const userData = {
                name: username,
                role: role
            };

            const response = await fetch('http://localhost:5554/signup', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(userData)
            });

            if (!response.ok) {
                throw new Error('Login failed');
            }

            setShowNotification('Login successful!');
            setShowNotification(true);
            // localStorage.setItem('access_token', );
            setUsername('');
            setRole('admin');
        } catch (error) {
            console.error('Error:', error);
            setShowNotification('Login failed. Please try again.');
            setShowNotification(true);
        }
    };

    const closeNotification = () => {
        setShowNotification(false);
    };

    return (
        <div className={"center-container"}>
            <div className={"container"} style={{ backgroundColor: 'darkcyan', padding: '20px', borderRadius: '20px' }}>
                <h1 style={{ textAlign: 'center', color: 'white' }}>Welcome to the Authorization Page</h1>
                <h2 style={{ textAlign: 'center'}}>First you need to sign up then login</h2>
                <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                    <div style={{ marginBottom: '20px' }}>
                        <label htmlFor="username" style={{ color: 'white', marginRight: '10px' }}>Username:</label>
                        <input
                            type="text"
                            id="username"
                            value={username}
                            placeholder="Enter your username"
                            onChange={(e) => setUsername(e.target.value)}
                            style={{ padding: '8px', borderRadius: '5px', border: '1px solid #ccc', marginRight: '10px' }}
                        />
                        <label htmlFor="role" style={{ color: 'white', marginRight: '10px' }}>Role:</label>
                        <select
                            id="role"
                            value={role}
                            onChange={(e) => setRole(e.target.value)}
                            style={{ padding: '8px', borderRadius: '5px', border: '1px solid #ccc', marginRight: '10px' }}
                        >
                            <option value="user">User</option>
                            <option value="admin">Admin</option>
                            {/* Add more options as needed */}
                        </select>
                    </div>
                    <div style={{ marginBottom: '10px' }}>
                        <button onClick={handleLogin} style={{ padding: '10px 20px', backgroundColor: 'green', color: 'white', border: 'none', borderRadius: '5px', cursor: 'pointer', marginRight: '10px' }}>Login</button>
                        <button onClick={handleSignup} style={{ padding: '10px 20px', backgroundColor: 'green', color: 'white', border: 'none', borderRadius: '5px', cursor: 'pointer' }}>Signup</button>
                        <button onClick={async () => {alert(getCookie("token"))}} style={{ padding: '10px 20px', margin: '10px', backgroundColor: 'green', color: 'white', border: 'none', borderRadius: '5px', cursor: 'pointer' }}>SHOW TOKEN</button>
                        <button onClick={() => window.location.href = '/tracker'} style={{ padding: '10px 20px', margin: '10px', backgroundColor: 'green', color: 'white', border: 'none', borderRadius: '5px', cursor: 'pointer' }}>To Service</button>
                        <button onClick={() => eraseCookie("token")} style={{ padding: '10px 20px', margin: '10px', backgroundColor: 'green', color: 'white', border: 'none', borderRadius: '5px', cursor: 'pointer' }}>Remove Token From Cookie</button>
                    </div>
                </div>
                {showNotification && <Notification message="Signup successful!" onClose={closeNotification} />}
            </div>
        </div>
    );
}

export default AuthPage;
