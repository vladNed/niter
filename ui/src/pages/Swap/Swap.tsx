import { AuthRedirectWrapper } from 'wrappers';
import { useScrollToElement } from 'hooks';
import { SwapWidget } from 'components';


export const Swap = () => {
  useScrollToElement();

  return (
    <AuthRedirectWrapper>
      <div className='h-full container flex flex-col gap-8 place-items-center'>
        <SwapWidget />
      </div>
    </AuthRedirectWrapper>
  );
};
