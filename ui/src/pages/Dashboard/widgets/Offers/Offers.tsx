import { Button, OutputContainer } from 'components';
import { useState } from 'react';
import { OfferData } from 'types';
import { CreateModal } from './components/CreateModal';

export const Offers = () => {
	const [isModalActive, setIsModalActive] = useState<boolean>(false);
	const [offers, setOffers] = useState<OfferData[]>([]);

	const toggleCreateModal = () => {
		setIsModalActive(!isModalActive);
	}

	const onSubmitOffer = async (data: OfferData) => {
		setOffers([...offers, data]);
		toggleCreateModal();
		try{
			const offerId = await wasmCreateOffer();
		} catch(e){
			console.error('Error creating offer', e);
		}

	}

	return (
		<OutputContainer>
			<div>
				<Button onClick={toggleCreateModal}>New Offer</Button>
				{isModalActive && <CreateModal onExit={toggleCreateModal} onSubmit={onSubmitOffer}/>}
			</div>
		</OutputContainer>
	);
}