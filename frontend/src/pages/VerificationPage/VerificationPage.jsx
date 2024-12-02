import React, { useState } from "react";
import Dialog from "../../components/UI/dialog/Dialog";
import "./style.css";

const VerificationPage = () => {
  const [otp, setOtp] = useState("");
  const [error, setError] = useState("");
  const [isVerified, setIsVerified] = useState(false);

  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const openDialog = () => {
    setIsDialogOpen(true);
  };

  const closeDialog = () => {
    setIsDialogOpen(false);
  };

  const handleOtpChange = (event) => {
    setOtp(event.target.value);
  };

  const handleVerify = () => {
    if (otp === "123456") {
      setIsVerified(true);
      setError("");
    } else {
      openDialog()
      setError("The entered code is incorrect.");
    }
  };

  return (
    <div className="verification-container">
      <h1>Identity Verification</h1>
      {isVerified ? (
        <div className="verification-success">
          <h2>Your account has been successfully verified!</h2>
        </div>
      ) : (
        <div className="verification-form">
          <label htmlFor="otp">Enter the verification code (OTP):</label>
          <input
            type="text"
            id="otp"
            value={otp}
            className={
              error ? "verification-input-error" : "verification-input"
            }
            onChange={handleOtpChange}
            placeholder="Verification code"
            maxLength="6"
            required
          />
          {error && <p className="error-message">{error}</p>}
          <button onClick={handleVerify}>Verify</button>
        </div>
      )}
      <Dialog isOpen={isDialogOpen} onClose={closeDialog} title="error">
        <p>The entered code is incorrect.</p>
        <button onClick={closeDialog}>Close</button>
      </Dialog>
    </div>
  );
};

export default VerificationPage;
