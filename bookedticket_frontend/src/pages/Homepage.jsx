import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { DndProvider } from 'react-dnd';
import { HTML5Backend } from 'react-dnd-html5-backend';
import KanbanColumn from '../components/KanbanColumn.jsx';
import '../styles/KanbanColumn.css';
import '../styles/KanbanCard.css';
import '../styles/Homepage.css';

const Homepage = () => {
    const [cards, setCards] = useState({
        pending: [],
        accepted: [],
        resolved: [],
        rejected: []
    });
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [editMode, setEditMode] = useState(null);
    const [editedData, setEditedData] = useState({ title: '', description: '', contact: '' });

    useEffect(() => {
        const fetchCards = async () => {
            try {
                const response = await axios.get('http://localhost:8080/api/v1/ticket');
                let data = response.data;

                if (!Array.isArray(data)) {
                    data = [data];
                }

                data.sort((a, b) => new Date(b.lastest_updated_at) - new Date(a.lastest_updated_at));

                setCards({
                    pending: data.filter(card => card.status === 'pending'),
                    accepted: data.filter(card => card.status === 'accepted'),
                    resolved: data.filter(card => card.status === 'resolved'),
                    rejected: data.filter(card => card.status === 'rejected')
                });
            } catch (error) {
                console.error('Error fetching cards:', error);
                setError('Failed to load tickets');
            } finally {
                setLoading(false);
            }
        };

        fetchCards();
    }, [cards]);

    const handleAddCard = async (newCard) => {
        try {
            const response = await axios.post('http://localhost:8080/api/v1/ticket', newCard, {
                headers: { "Content-Type": "application/json" }
            });

            setCards(prevCards => ({
                ...prevCards,
                [newCard.status]: [...prevCards[newCard.status], response.data]
            }));
        } catch (error) {
            console.error('Error adding card:', error);
        }
    };

    const handleSave = async (cardId, status) => {
        try {
            await axios.put(`http://localhost:8080/api/v1/ticket/${cardId}`, {
                title: editedData.title,
                description: editedData.description,
                status: status,
                contact: editedData.contact
            }, {
                headers: { "Content-Type": "application/json" }
            });

            setCards(prevCards => ({
                ...prevCards,
                [status]: prevCards[status].map(card =>
                    card.id === cardId ? { ...card, ...editedData } : card
                )
            }));
            setEditMode(null);
        } catch (error) {
            console.error('Error updating card:', error);
            setError('Failed to update ticket');
        }
    };

    const handleCardUpdate = (cardId, updatedData) => {
        setCards(prevCards => {
            const updatedCards = { ...prevCards };

            Object.keys(updatedCards).forEach(status => {
                updatedCards[status] = updatedCards[status].map(card =>
                    card.id === cardId ? { ...card, ...updatedData } : card
                );
            });

            return updatedCards;
        });
    };


    const handleDrop = async (cardId, newStatus) => {
        try {
            const cardToUpdate = Object.values(cards).flat().find(card => card.id === cardId);

            if (!cardToUpdate) {
                console.error("Card not found");
                return;
            }

            await axios.put(`http://localhost:8080/api/v1/ticket/${cardId}`, {
                title: cardToUpdate.title,
                description: cardToUpdate.description,
                contact: cardToUpdate.contact,
                status: newStatus
            });

            setCards(prevCards => {
                let updatedCards = { ...prevCards };

                let movedCard;
                Object.keys(updatedCards).forEach(status => {
                    updatedCards[status] = updatedCards[status].filter(card => {
                        if (card.id === cardId) {
                            movedCard = { ...card, status: newStatus };
                            return false;
                        }
                        return true;
                    });
                });

                if (movedCard) {
                    updatedCards[newStatus] = [...updatedCards[newStatus], movedCard];
                }

                return updatedCards;
            });

        } catch (error) {
            console.error("Error updating card status:", error);
        }
    };


    return (
        <DndProvider backend={HTML5Backend}>
            <div className='container'>
                {loading && <div>Loading...</div>}  
                {error && <div>{error}</div>}
                <div className="kanban-board">
                    {['pending', 'accepted', 'resolved', 'rejected'].map(status => (
                        <KanbanColumn
                            key={status}
                            status={status}
                            cards={cards[status]}
                            onDrop={handleDrop}
                            onCardUpdate={handleCardUpdate}
                            onAddCard={handleAddCard}
                        />
                    ))}
                </div>

                {editMode && (
                    <div className="edit-modal">
                        <h2>Edit Ticket</h2>
                        <input
                            type="text"
                            value={editedData.title}
                            onChange={(e) => setEditedData({ ...editedData, title: e.target.value })}
                        />
                        <textarea
                            value={editedData.description}
                            onChange={(e) => setEditedData({ ...editedData, description: e.target.value })}
                        />
                        <input
                            type="text"
                            value={editedData.contact}
                            onChange={(e) => setEditedData({ ...editedData, contact: e.target.value })}
                        />
                        <button onClick={() => {
                            const foundCard = Object.values(cards).flat().find(card => card.id === editMode);
                            if (foundCard) handleSave(editMode, foundCard.status);
                        }}>Save</button>
                        <button onClick={() => setEditMode(null)}>Cancel</button>
                    </div>
                )}
            </div>
        </DndProvider>
    );
};

export default Homepage;