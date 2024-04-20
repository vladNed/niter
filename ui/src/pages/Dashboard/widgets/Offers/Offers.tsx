import { Button, OutputContainer } from 'components';
import { useEffect, useState } from 'react';
import { OfferData } from 'types';
import { CreateModal } from './components/CreateModal';
import { SIGNALLING_SERVER_URL } from 'config';

export const Offers = () => {
	const [isCreateModalActive, setIsCreateModalActive] = useState<boolean>(false);
	const [isAnswerModalActive, setIsAnswerModalActive] = useState<boolean>(false);
	const [offers, setOffers] = useState<string[]>([]);

	const toggleCreateModal = () => {
		setIsCreateModalActive(!isCreateModalActive);
	}

	const toggleAnswerModal = () => {
		setIsAnswerModalActive(!isAnswerModalActive);
	}

	useEffect(() => {
		const ws = new WebSocket(SIGNALLING_SERVER_URL);

		ws.onopen = () => {
			console.log('Connected to signalling server')
			const registration = JSON.stringify({ channel: ['marketplace']});
			ws.send(registration);
		}

		ws.onmessage = (event: MessageEvent<any>) => {
			const message = JSON.parse(event.data);
			console.log(message)
			if (message.type === 'offer') {
				setOffers((offers) => [...offers, message.offerId]);
			}
		}



		return () => {
			ws.close();
		}

	}, []);

	const onSubmitOffer = async (data: OfferData) => {
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
				<Button onClick={toggleAnswerModal}>Connect Offer</Button>
				{offers.map((offerId) => (
					<div key={offerId}>
						{offerId}
					</div>
				))}
				{isCreateModalActive && <CreateModal onExit={toggleCreateModal} onSubmit={onSubmitOffer}/>}
			</div>
		</OutputContainer>
	);
}

