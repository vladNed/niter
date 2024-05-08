import { AuthRedirectWrapper } from 'wrappers';
import { SearchOfferWidget } from 'components';

export const Home = () => {
  return (
    <AuthRedirectWrapper requireAuth={false}>
      <SearchOfferWidget isPlaceholder={true}/>
    </AuthRedirectWrapper>
  );
};
