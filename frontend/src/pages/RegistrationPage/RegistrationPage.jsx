import { useRef, useState } from "react";
import Dialog from "../../components/UI/dialog/Dialog";
import "./style.css";

const RegistrationPage = () => {
  const [firstNameError, setFirstNameError] = useState("");
  const [lastNameError, setLastNameError] = useState("");
  const [phoneError, setPhoneError] = useState("");
  const [emailError, setEmailError] = useState("");
  const [passwordError, setPasswordError] = useState("");

  const inputFirstName = useRef();
  const inputLastName = useRef();
  const inputPhone = useRef();
  const inputEmail = useRef();
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
      firstname: inputFirstName.current.value,
      lastname: inputLastName.current.value,
      phone: inputPhone.current.value,
      email: inputEmail.current.value,
      password: inputPassword.current.value,
    };

    let hasError = false;

    // Validation for First Name
    if (!formData.firstname.trim()) {
      setFirstNameError("First name is required");
      hasError = true;
    } else {
      setFirstNameError("");
    }

    if (!formData.lastname.trim()) {
      setLastNameError("Last name is required");
      hasError = true;
    } else {
      setLastNameError("");
    }

    // Validation for Phone
    const phoneRegex = /^[0-9]{10}$/;
    if (!phoneRegex.test(formData.phone)) {
      setPhoneError("Enter a valid 10-digit phone number");
      hasError = true;
    } else {
      setPhoneError("");
    }

    // Validation for Email
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(formData.email)) {
      setEmailError("Enter a valid email address");
      hasError = true;
    } else {
      setEmailError("");
    }

    // Validation for Password
    if (formData.password.length < 6) {
      setPasswordError("Password must be at least 6 characters");
      hasError = true;
    } else {
      setPasswordError("");
    }

    // If no errors, handle form submission
    if (!hasError) {
      console.log("Form submitted successfully");
    } else {
      openDialog();
    }
  };

  return (
    <div className="registration">
      <h1>Welcome, Create Your Account</h1>
      <form className="registration-form" onSubmit={handleSubmit}>
        <label>
          First Name:
          <input
            ref={inputFirstName}
            className={
              firstNameError ? "registration-input-error" : "registration-input"
            }
            type="text"
            placeholder="Enter your first name"
          />
        </label>
        <label>
          Last Name:
          <input
            ref={inputLastName}
            className={
              lastNameError ? "registration-input-error" : "registration-input"
            }
            type="text"
            placeholder="Enter your last name"
          />
        </label>
        <label>
          Phone:
          <input
            ref={inputPhone}
            className={
              phoneError ? "registration-input-error" : "registration-input"
            }
            type="tel"
            placeholder="Enter your phone number"
          />
        </label>
        <label>
          Email:
          <input
            ref={inputEmail}
            className={
              emailError ? "registration-input-error" : "registration-input"
            }
            type="email"
            placeholder="Enter your email"
          />
        </label>
        <label>
          Password:
          <input
            ref={inputPassword}
            className={
              passwordError ? "registration-input-error" : "registration-input"
            }
            type="password"
            placeholder="Enter your password"
          />
        </label>
        <button className="submit-button" type="submit">Sign Up</button>
      </form>
      <Dialog isOpen={isDialogOpen} onClose={closeDialog} title="error">
        <p>fill your profile information.</p>
        <button onClick={closeDialog}>Close</button>
      </Dialog>
    </div>
  );
};

export default RegistrationPage;
