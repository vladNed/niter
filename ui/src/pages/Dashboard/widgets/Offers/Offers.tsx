import { Button, OutputContainer } from 'components';
import { useEffect, useState } from 'react';
import { OfferData } from 'types';
import { CreateModal } from './components/CreateModal';
import { OFFERS_POLLING_INTERVAL } from 'config';

export const Offers = () => {
	const [isCreateModalActive, setIsCreateModalActive] = useState<boolean>(false);
	const [offers, setOffers] = useState<string[]>([]);

	const toggleCreateModal = () => {
		setIsCreateModalActive(!isCreateModalActive);
	}

	const onOffersPool = async () => {
		try {
			const newOffers = await wasmPollOffers();
			setOffers(newOffers.sort());
		} catch (e) {
			console.error('Error polling offers', e);
		}
	}

	useEffect(() => {
		setInterval(onOffersPool, OFFERS_POLLING_INTERVAL);
	}, [onOffersPool]);

	const onSubmitOffer = async (data: OfferData) => {
		toggleCreateModal();
		try {
			await wasmCreateOffer();
		} catch (e) {
			console.error('Error creating offer', e);
		}
	}

	const onSubmitAnswer = async (offerId: string) => {
		try {
			await wasmCreateAnswer(offerId);
		} catch (e) {
			console.error('Error creating offer', e);
		}
	}

	return (
		<OutputContainer>
			<div>
				<Button onClick={toggleCreateModal}>New Offer</Button>
				<Button onClick={onOffersPool}>Refresh</Button>
				{offers.map((offerId) => (
					<div key={offerId} className='flex items-center gap-2 mt-2 border-[1px] rounded mt-2 p-2 hover:shadow-xl duration-200 ease-in hover:bg-sky-100'>
						<div>Offer ID: {offerId}</div>
						<Button onClick={() => onSubmitAnswer(offerId)}>Connect</Button>
					</div>
				))}
				{isCreateModalActive && <CreateModal onExit={toggleCreateModal} onSubmit={onSubmitOffer} />}
			</div>
		</OutputContainer>
	);
}

