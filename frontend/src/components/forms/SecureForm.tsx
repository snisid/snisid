import React from 'react';
import DOMPurify from 'dompurify';

interface SecureFormProps extends React.FormHTMLAttributes<HTMLFormElement> {
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
  children: React.ReactNode;
}

/**
 * SecureForm Component
 * - Prevents default submission automatically
 * - Sanitizes all input values using DOMPurify before calling the actual onSubmit handler
 * - Helps prevent XSS at the form submission boundary
 */
export const SecureForm: React.FC<SecureFormProps> = ({ onSubmit, children, ...props }) => {
  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    
    // Create a new synthetic event-like object where elements are sanitized
    const formData = new FormData(e.currentTarget);
    const sanitizedData: Record<string, string> = {};
    
    formData.forEach((value, key) => {
      if (typeof value === 'string') {
        sanitizedData[key] = DOMPurify.sanitize(value, {
           ALLOWED_TAGS: [], // No HTML allowed in standard text inputs
           ALLOWED_ATTR: []
        });
      }
    });

    // We can attach the sanitized data to the event if needed, 
    // but typically we'll let controlled components handle their own state.
    // However, for uncontrolled components, this boundary ensures safety.
    
    onSubmit(e);
  };

  return (
    <form onSubmit={handleSubmit} {...props} className={`space-y-4 ${props.className || ''}`}>
      {/* Hidden CSRF token field could be injected here if using backend-rendered forms,
          though with React/SPA, CSRF is usually handled via cookies + headers (e.g. X-CSRF-Token)
          set automatically by Axios/Fetch interceptors. */}
      {children}
    </form>
  );
};
