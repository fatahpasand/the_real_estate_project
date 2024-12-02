import { useRef, useState } from "react";
import Dialog from "../../components/UI/dialog/Dialog";
import "./style.css";

const LoginPage = () => {
  const [phoneError, setPhoneError] = useState("");
  const [passwordError, setPasswordError] = useState("");

  const inputPhone = useRef();
  const inputPassword = useRef();

  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const openDialog = () => {
    setIsDialogOpen(true);
  };

  const closeDialog = () => {
    setIsDialogOpen(false);
  };
  const handleSubmit = (e) => {
    e.preventDefault();
    const formData = {
      phone: inputPhone.current.value,
      password: inputPassword.current.value,
    };

    let hasError = false;

    const phoneRegex = /^[0-9]{10}$/;
    if (!phoneRegex.test(formData.phone)) {
      setPhoneError("Enter a valid 10-digit phone number");
      hasError = true;
    } else {
      setPhoneError("");
    }

    if (formData.password.length < 6) {
      setPasswordError("Password must be at least 6 characters");
      hasError = true;
    } else {
      setPasswordError("");
    }

    if(hasError){
      openDialog()
    }

  
  };

  return (
    <div className="login">
      <h1>Welcome Back</h1>
      <form className="login-form" onSubmit={handleSubmit}>
        <label>
          Phone:
          <input
            ref={inputPhone}
            className={phoneError ? "login-input-error" : "login-input"}
            type="tel"
            placeholder="Enter your phone number"
          />
        </label>
        <label>
          Password:
          <input
            ref={inputPassword}
            className={passwordError ? "login-input-error" : "login-input"}
            type="password"
            placeholder="Enter your password"
          />
        </label>
        <button type="submit">Sign Up</button>
      </form>
      <Dialog isOpen={isDialogOpen} onClose={closeDialog} title="error">
        <p>fill your phone and password.</p>
        <button onClick={closeDialog}>Close</button>
      </Dialog>
    </div>
  );
};

export default LoginPage;
