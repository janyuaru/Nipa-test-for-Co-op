import React, { useState } from 'react';
import { useDrop } from 'react-dnd';
import KanbanCard from './KanbanCard';
import '../styles/KanbanColumn.css';

const KanbanColumn = ({ status, cards, onDrop, onCardUpdate, onAddCard }) => {
    const [, drop] = useDrop({
        accept: 'CARD',
        drop: (item) => onDrop(item.id, status),
    });

    const [showForm, setShowForm] = useState(false);
    const [newCardData, setNewCardData] = useState({ title: '', description: '', contact: '', status: status });

    const handleAdd = async () => {
        if (newCardData.title.trim()) {
            onAddCard(newCardData);
            setNewCardData({ title: '', description: '', contact: '', status: status });
            setShowForm(false);
        }
    };


    return (
        <div ref={drop} className="kanban-column">
            <h2>{status.charAt(0).toUpperCase() + status.slice(1)}</h2>

            <button className="kanban-add-btn" onClick={() => setShowForm(true)}>+ Add Task</button>
            {showForm && (
                <div className="add-form">
                    <input
                        type="text"
                        placeholder="Title"
                        value={newCardData.title}
                        onChange={(e) => setNewCardData({ ...newCardData, title: e.target.value })}
                    />
                    <textarea
                        placeholder="Description"
                        value={newCardData.description}
                        onChange={(e) => setNewCardData({ ...newCardData, description: e.target.value })}
                    />
                    <input
                        type="text"
                        placeholder="Contact"
                        value={newCardData.contact}
                        onChange={(e) => setNewCardData({ ...newCardData, contact: e.target.value })}
                    />
                    <button className='button-save' onClick={handleAdd}>Save</button>
                    <button className='button-cancel' onClick={() => setShowForm(false)}>Cancel</button>
                </div>
            )}

            <div className="kanban-card-container">
                {cards.map((card) => (
                    <KanbanCard
                        key={card.id}
                        card={card}
                        onCardUpdate={onCardUpdate}
                    />
                ))}
            </div>

        </div>
    );
};

export default KanbanColumn;