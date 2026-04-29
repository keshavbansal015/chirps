import React, {useState} from 'react';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';

import Link from '../../router/link';
import UserDropdown from './dropdown';
import {useTheme} from '../../contexts/ThemeContext';
import './navbar.scss';

function Navbar(props) {
  const [state, setState] = useState({
    isDropdownShown: false
  });
  const {isDarkMode, toggleDarkMode} = useTheme();

  function handleClick() {
    setState((state) => ({
      ...state,
      isDropdownShown: !state.isDropdownShown
    }));
  }

  function handleThemeToggle(e) {
    e.stopPropagation();
    toggleDarkMode();
  }

  return (
    <header className='navigation-container bottom-shadow'>
      <div className='navigation-content main-container'>
        {props.isLoggedIn ? (
          <div className='left-items'>
            <Link href='/feed' className='nav-button'>
              <FontAwesomeIcon
                icon='home'
                className='nav-button-icon'
                size='2x'
              />
              <span className='nav-button-label'>Home</span>
            </Link>
          </div>
        ) : (
          <div className='left-items'></div>
        )}

        <FontAwesomeIcon icon='comment-dots' className='icon' size='2x' />

        {props.user ? (
          <div className='right-items'>
            <button
              className='theme-toggle'
              onClick={handleThemeToggle}
              aria-label={isDarkMode ? 'Switch to light mode' : 'Switch to dark mode'}
            >
              <FontAwesomeIcon
                icon={isDarkMode ? 'sun' : 'moon'}
                className='theme-toggle-icon'
              />
            </button>
            <div className='profile-button' onClick={handleClick}>
              <img
                className='profile-button-image'
                src='https://via.placeholder.com/300.png'
                alt='Profile'
              />

              {state.isDropdownShown &&
                <UserDropdown
                  user={props.user}
                  logoutUser={props.logoutUser}
                />
              }
            </div>
          </div>
        ) : (
          <div className='right-items'>
            <button
              className='theme-toggle'
              onClick={handleThemeToggle}
              aria-label={isDarkMode ? 'Switch to light mode' : 'Switch to dark mode'}
            >
              <FontAwesomeIcon
                icon={isDarkMode ? 'sun' : 'moon'}
                className='theme-toggle-icon'
              />
            </button>
            <Link href='/login' className='button login-button'>
              Log In
            </Link>
          </div>
        )}
      </div>
    </header>
  );
}

export default Navbar;
