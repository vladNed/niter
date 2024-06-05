import { Button } from 'components/Button';
import { MxLink } from 'components/MxLink';
import { environment } from 'config';
import { logout } from 'helpers';
import { useGetIsLoggedIn } from 'hooks';
import { RouteNamesEnum } from 'localConstants';
import { useMatch } from 'react-router-dom';
import { Logo, NavigationItems } from 'components/Layout/Header/Components';
import { NetworkDrawer, NetworkModal } from 'components/Modal';
import { useState } from 'react';

const callbackUrl = `${window.location.origin}/unlock`;
const onRedirect = undefined; // use this to redirect with useNavigate to a specific page after logout
const shouldAttemptReLogin = false; // use for special cases where you want to re-login after logout
const options = {
  /*
   * @param {boolean} [shouldBroadcastLogoutAcrossTabs=true]
   * @description If your dApp supports multiple accounts on multiple tabs,
   * this param will broadcast the logout event across all tabs.
   */
  shouldBroadcastLogoutAcrossTabs: true,
  /*
   * @param {boolean} [hasConsentPopup=false]
   * @description Set it to true if you want to perform async calls before logging out on Safari.
   * It will open a consent popup for the user to confirm the action before leaving the page.
   */
  hasConsentPopup: false
};

export const Header = () => {
  const isLoggedIn = useGetIsLoggedIn();
  const isUnlockRoute = Boolean(useMatch(RouteNamesEnum.unlock));
  const [drawerOpen, setDrawerOpen] = useState<boolean>(false);

  const handleDrawer = () => {
    setDrawerOpen(!drawerOpen);
  }

  const ConnectButton = isUnlockRoute ? null : (
    <MxLink to={RouteNamesEnum.unlock}>Connect</MxLink>
  );

  const handleLogout = () => {
    sessionStorage.clear();
    logout(
      callbackUrl,
      /*
       * following are optional params. Feel free to remove them in your implementation
       */
      onRedirect,
      shouldAttemptReLogin,
      options
    );
  };

  const navItems = [
    { title: 'Swap', path: RouteNamesEnum.swap },
  ]

  return (
    <header className='flex flex-row align-center gap-6 pl-6 pr-6 pt-6 pb-6 w-full m-auto justify-between'>
      <Logo />
      <nav className='float-end text-white h-full w-full text-sm sm:relative sm:left-auto sm:top-auto sm:flex sm:w-auto sm:flex-row sm:justify-end sm:bg-transparent p-1'>
        <div className='flex float-end container mx-auto items-center gap-6'>
          <div className='flex gap-1 items-center text-black'>
            <div className='w-2 h-2 rounded-full bg-green-600' />
            <p>{environment}</p>
          </div>
          <button disabled={!isLoggedIn} onClick={handleDrawer} className='px-3 py-2 rounded-lg bg-primary-500 hover:no-underline hover:bg-primary-600/100 transition-colors hover:text-primary-100'>
            Network
          </button>

          {isLoggedIn ? (
            <Button
              onClick={handleLogout}
              className='text-black inline-block rounded-lg px-3 py-2 text-center hover:no-underline my-0 hover:bg-slate-300 mx-0 bg-slate-200 '
            >
              Close
            </Button>
          ) : (
            ConnectButton
          )}
        </div>
      <NetworkDrawer isOpen={drawerOpen} toggleDrawer={handleDrawer} />
      </nav>
    </header>
  );
};
