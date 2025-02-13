import React, { useState } from 'react';
import { useDrag } from 'react-dnd';
import axios from 'axios';
import '../styles/KanbanCard.css';

const KanbanCard = ({ card, onCardUpdate }) => {
    const [{ isDragging }, drag] = useDrag({
        type: 'CARD',
        item: { id: card.id, status: card.status },
        collect: (monitor) => ({
            isDragging: monitor.isDragging(),
        }),
    });

    const [isEditing, setIsEditing] = useState(false);
    const [editedData, setEditedData] = useState({
        title: card.title,
        description: card.description,
        contact: card.contact
    });

    const handleSave = async () => {
        try {
            await axios.put(`http://localhost:8080/api/v1/ticket/${card.id}`, {
                ...editedData,
                status: card.status,
            });

            onCardUpdate(card.id, editedData);
            setIsEditing(false);
        } catch (error) {
            console.error("Error updating card:", error);
        }
    };

    return (
        <div ref={drag} className="kanban-card" style={{ opacity: isDragging ? 0.5 : 1 }}>
            {isEditing ? (
                <>
                    <input
                        type="text"
                        className='kanban-input'
                        placeholder="Title"
                        value={editedData.title}
                        onChange={(e) => setEditedData({ ...editedData, title: e.target.value })}
                    />
                    <textarea
                        className='kanban-textarea'
                        placeholder="Description"
                        value={editedData.description}
                        onChange={(e) => setEditedData({ ...editedData, description: e.target.value })}
                    />
                    <input
                        type="text"
                        className='kanban-input'
                        placeholder="Contact"
                        value={editedData.contact}
                        onChange={(e) => setEditedData({ ...editedData, contact: e.target.value })}
                    />
                    <div className="kanban-card-edit">
                        <button className="kanban-btn save" onClick={handleSave}>Save</button>
                        <button className="kanban-btn cancel" onClick={() => setIsEditing(false)}>Cancel</button>
                    </div>
                </>
            ) : (
                <>
                    <h3>{card.title}</h3>
                    <p>{card.description}</p>
                    <p>{card.contact}</p>
                    <div className="kanban-card-edit">
                        <button
                            className='kanban-btn'
                            onClick={(e) => {
                                e.stopPropagation();
                                setIsEditing(true);
                            }}>Edit</button>
                    </div>
                </>
            )}
        </div>
    );
};


export default KanbanCard;