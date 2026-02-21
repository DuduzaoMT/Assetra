// Input validation and sanitization utilities
// Helps prevent XSS and injection attacks

/**
 * Validates email format
 * @param {string} email - Email to validate
 * @returns {boolean} - True if valid
 */
export const isValidEmail = (email) => {
  const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  return emailRegex.test(email) && email.length <= 100;
};

/**
 * Validates password strength
 * @param {string} password - Password to validate
 * @returns {object} - { isValid: boolean, errors: string[] }
 */
export const validatePassword = (password) => {
  const errors = [];

  if (!password) {
    errors.push("Password is required");
    return { isValid: false, errors };
  }

  if (password.length < 8) {
    errors.push("Password must be at least 8 characters");
  }

  if (password.length > 128) {
    errors.push("Password is too long");
  }

  if (!/[A-Z]/.test(password)) {
    errors.push("Password must contain at least one uppercase letter");
  }

  if (!/[a-z]/.test(password)) {
    errors.push("Password must contain at least one lowercase letter");
  }

  if (!/[0-9]/.test(password)) {
    errors.push("Password must contain at least one number");
  }

  // eslint-disable-next-line no-useless-escape
  if (!/[!@#$%^&*()_+=\[\]{};':"\\|,.<>/?-]/.test(password)) {
    errors.push("Password must contain at least one special character");
  }

  return {
    isValid: errors.length === 0,
    errors,
  };
};

/**
 * Validates username/name
 * @param {string} name - Name to validate
 * @returns {object} - { isValid: boolean, error: string }
 */
export const validateName = (name) => {
  if (!name || name.trim().length === 0) {
    return { isValid: false, error: "Name is required" };
  }

  const trimmed = name.trim();

  if (trimmed.length < 2) {
    return { isValid: false, error: "Name must be at least 2 characters" };
  }

  if (trimmed.length > 50) {
    return { isValid: false, error: "Name must be at most 50 characters" };
  }

  // Check for suspicious characters
  if (/[<>&"']/.test(trimmed)) {
    return { isValid: false, error: "Name contains invalid characters" };
  }

  return { isValid: true, error: null };
};

/**
 * Sanitizes string input to prevent XSS
 * @param {string} input - Input to sanitize
 * @returns {string} - Sanitized input
 */
export const sanitizeInput = (input) => {
  if (!input) return "";

  return input
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#x27;")
    .replace(/\//g, "&#x2F;");
};

/**
 * Sanitizes name by removing dangerous characters
 * @param {string} name - Name to sanitize
 * @returns {string} - Sanitized name
 */
export const sanitizeName = (name) => {
  if (!name) return "";

  return name.trim().replace(/[<>&"']/g, "");
};

/**
 * Validates and sanitizes form data
 * @param {object} formData - Form data to validate
 * @param {object} schema - Validation schema
 * @returns {object} - { isValid: boolean, errors: object, sanitized: object }
 */
export const validateFormData = (formData, schema) => {
  const errors = {};
  const sanitized = {};
  let isValid = true;

  Object.keys(schema).forEach((field) => {
    const value = formData[field];
    const rules = schema[field];

    if (rules.required && (!value || value.trim() === "")) {
      errors[field] = `${field} is required`;
      isValid = false;
      return;
    }

    if (rules.type === "email" && value) {
      if (!isValidEmail(value)) {
        errors[field] = "Invalid email format";
        isValid = false;
        return;
      }
      sanitized[field] = value.trim().toLowerCase();
    }

    if (rules.type === "password" && value) {
      const validation = validatePassword(value);
      if (!validation.isValid) {
        errors[field] = validation.errors.join(". ");
        isValid = false;
        return;
      }
      sanitized[field] = value; // Don't trim passwords
    }

    if (rules.type === "text" && value) {
      const nameValidation = validateName(value);
      if (!nameValidation.isValid) {
        errors[field] = nameValidation.error;
        isValid = false;
        return;
      }
      sanitized[field] = sanitizeName(value);
    }
  });

  return { isValid, errors, sanitized };
};
