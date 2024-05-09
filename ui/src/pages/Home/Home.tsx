import { AuthRedirectWrapper } from 'wrappers';
import { SwapWidget } from 'components';

export const Home = () => {
  return (
    <AuthRedirectWrapper requireAuth={false}>
      <SwapWidget isPlaceholder={true}/>
    </AuthRedirectWrapper>
  );
};
