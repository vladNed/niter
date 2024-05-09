import { AuthRedirectWrapper } from 'wrappers';
import { useScrollToElement } from 'hooks';
import { LoadWasm } from 'wasm';
import { LoadingModal, SearchOfferWidget } from 'components';


export const Dashboard = () => {
  useScrollToElement();

  return (
    <AuthRedirectWrapper>
      <LoadWasm>
        <div className='h-full container flex flex-col gap-8 place-items-center'>
          <SearchOfferWidget isPlaceholder={false} />
        </div>
      </LoadWasm>
    </AuthRedirectWrapper>
  );
};
