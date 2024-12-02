import React, { useState } from 'react'

import styles from "./styles.module.scss"
import { Link } from 'react-router-dom'
import Dialog from '../../UI/dialog/Dialog'
import RegistrationPage from '../../../pages/RegistrationPage/RegistrationPage'
import LoginPage from '../../../pages/LoginPage/LoginPage'


export default function Header() {

    const [showRegister, setShowRegister] = useState(false)

    const [showLogin, setShowLogin] = useState(false)


    return (
        <header className={styles.header}>
            <div className={styles.left}>
                <Link to={"/"}>
                    <img src='/images/logo.png' />
                    <span>Estate</span>
                </Link>
            </div>
            <div className={styles.right}>
                <h1 onClick={() => {
                    setShowRegister(prev => !prev)
                }}>
                    Sign up</h1>
                <h1 onClick={() => {
                    setShowLogin(prev => !prev)
                }}>
                    Login</h1>
            </div>

            <Dialog isOpen={showRegister} onClose={() => setShowRegister(false)}
                title={"Register"}>
                <RegistrationPage />
            </Dialog>
            <Dialog isOpen={showLogin} onClose={() => setShowLogin(false)}
                title={"Login"}>
                <LoginPage />
            </Dialog>
        </header>
    )
}
