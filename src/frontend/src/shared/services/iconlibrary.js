import {library} from '@fortawesome/fontawesome-svg-core';

import {
  faCommentDots, faHome, faTimes, faRetweet,
  faHeart, faCalendarAlt, faPassport,
  faUser,  faEnvelope,  faLock,
  faEye, faMobileAlt, faDesktop,
  faFileAlt, faImage, faSun, faMoon
} from '@fortawesome/free-solid-svg-icons';

const IconLibrary = {
  configure: () => {
    library.add(
      faCommentDots, faHome, faTimes, faRetweet,
      faHeart, faCalendarAlt, faPassport,
      faUser, faEnvelope, faLock,
      faEye, faMobileAlt, faDesktop,
      faFileAlt, faImage, faSun, faMoon
    );
  }
};

export default IconLibrary;
