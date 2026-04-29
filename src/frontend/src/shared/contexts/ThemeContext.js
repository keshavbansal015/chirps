import React, {createContext, useContext, useState, useEffect} from 'react';

const ThemeContext = createContext();

export function ThemeProvider({children}) {
  const [isDarkMode, setIsDarkMode] = useState(() => {
    const saved = localStorage.getItem('darkMode');
    return saved ? JSON.parse(saved) : false;
  });

  useEffect(() => {
    localStorage.setItem('darkMode', JSON.stringify(isDarkMode));
    document.documentElement.setAttribute('data-theme', isDarkMode ? 'dark' : 'light');
  }, [isDarkMode]);

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', isDarkMode ? 'dark' : 'light');
  }, []);

  const toggleDarkMode = () => setIsDarkMode(prev => !prev);

  return (
    <ThemeContext.Provider value={{isDarkMode, toggleDarkMode}}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return context;
}
