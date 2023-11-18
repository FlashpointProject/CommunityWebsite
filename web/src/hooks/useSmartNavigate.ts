import { NavigateFunction, NavigateOptions, To, useNavigate } from 'react-router-dom';

/**
 * Modifies the useNavigate hook to prevent navigation to the same page you're on, unless a query string is provided.
 */
export function useSmartNavigate(): NavigateFunction {
  const navigate = useNavigate();
  const customNavigate: NavigateFunction = (toOrDelta: To | number, options?: NavigateOptions): void => {
    if (typeof toOrDelta === 'number') {
      // Delta navigation
      navigate(toOrDelta);
    } else {
      // Path / string navigation
      console.log(toOrDelta);
      if (toOrDelta === window.location.pathname) {
        return;
      } else {
        navigate(toOrDelta, options);
      }
      navigate(toOrDelta, options);
    }
  };
  return customNavigate;
}
