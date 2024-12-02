import { useState, useRef } from "react";
import Dialog from "../../components/UI/dialog/Dialog";
import "./style.css";
const DashboardPage = () => {
  const [firstNameError, setFirstNameError] = useState("");
  const [lastNameError, setLastNameError] = useState("");
  const [idError, setIdError] = useState("");
  const [ageError, setAgeError] = useState("");
  const [phoneError, setPhoneError] = useState("");
  const [emailError, setEmailError] = useState("");

  const inputFirstName = useRef();
  const inputLastName = useRef();
  const inputID = useRef();
  const inputAge = useRef();
  const inputPhone = useRef();
  const inputEmail = useRef();

  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const openDialog = () => {
    setIsDialogOpen(true);
  };

  const closeDialog = () => {
    setIsDialogOpen(false);
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    let hasError = false;

    // Validation for First Name
    if (!inputFirstName.current.value) {
      setFirstNameError("First name is required");
      hasError = true;
    } else {
      setFirstNameError("");
    }

    // Validation for Last Name
    if (!inputLastName.current.value) {
      setLastNameError("Last name is required");
      hasError = true;
    } else {
      setLastNameError("");
    }

    // Validation for ID
    if (!inputID.current.value) {
      setIdError("Id is required");
      hasError = true;
    } else {
      setIdError("");
    }
    // Validation for Age
    const ageRegex = /^[0-9]/;
    if (!ageRegex.test(inputAge.current.value)) {
      setAgeError("Enter a age number");
      hasError = true;
    } else {
      setAgeError("");
    }

    // Validation for Phone
    const phoneRegex = /^[0-9]{10}$/;
    if (!phoneRegex.test(inputPhone.current.value)) {
      setPhoneError("Enter a valid 10-digit phone number");
      hasError = true;
    } else {
      setPhoneError("");
    }

    // Validation for Email
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(inputEmail.current.value)) {
      setEmailError("Enter a valid email address");
      hasError = true;
    } else {
      setEmailError("");
    }

    // If no errors, handle form submission
    if (!hasError) {
      const formData = {
        fistname: inputFirstName.current.value,
        lastname: inputLastName.current.value,
        id: inputID.current.value,
        age: inputAge.current.value,
        phone: inputPhone.current.value,
        email: inputEmail.current.value,
      };

      console.log("Form submitted successfully");
    }else {
      openDialog();
    }
  };
  return (
    <div className="dashboard">
      <div className="dashboard-header">
        <h1>Dashboard</h1>
      </div>
      <div className="dashboard-content">
        <form className="dashboard-form">
          <label>
            First Name:
            <input
              type="text"
              ref={inputFirstName}
              className={
                firstNameError ? "dashboard-input-error" : "dashboard-input"
              }
              placeholder="Enter your fist name"
            ></input>
          </label>
          <label>
            Last Name:
            <input
              type="text"
              ref={inputLastName}
              className={
                lastNameError ? "dashboard-input-error" : "dashboard-input"
              }
              placeholder="Enter your last name"
            ></input>
          </label>
          <label>
            ID:
            <input
              type="text"
              ref={inputID}
              className={idError ? "dashboard-input-error" : "dashboard-input"}
              placeholder="Enter your id"
            ></input>
          </label>
          <label>
            Age:
            <input
              type="number"
              min="0"
              max="200"
              ref={inputAge}
              className={ageError ? "dashboard-input-error" : "dashboard-input"}
              placeholder="Enter your age"
            ></input>
          </label>
          <label>
            Phone:
            <input
              type="text"
              ref={inputPhone}
              className={
                phoneError ? "dashboard-input-error" : "dashboard-input"
              }
              placeholder="Enter your phone"
            ></input>
          </label>
          <label>
            Email:
            <input
              type="text"
              ref={inputEmail}
              className={
                emailError ? "dashboard-input-error" : "dashboard-input"
              }
              placeholder="Enter your email"
            ></input>
          </label>
        </form>

        <div className="dashboard-edits">
          <div>
            <i className=""></i>
            <button>Edit Profile</button>
          </div>
          <div>
            <i className=""></i>
            <button>Deals History</button>
          </div>
          <div>
            <i className=""></i>
            <button>Favorites</button>
          </div>
          <div>
            <i className=""></i>
            <button>Register ad</button>
          </div>
        </div>
      </div>
      <div className="dashboard-buttons">
        <button>Delete Account</button>
        <button onClick={handleSubmit}>Done</button>
      </div>
      <div className="dashboard-footer"></div>

      <Dialog isOpen={isDialogOpen} onClose={closeDialog} title="error">
        <p>fill your profile information.</p>
        <button onClick={closeDialog}>Close</button>
      </Dialog>
    </div>
  );
};

export default DashboardPage;
