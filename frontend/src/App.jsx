import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';

function App() {
    const [names, setNames] = useState([]);  // Initialize as an empty array
    const [newName, setNewName] = useState('');
    const [updateName, setUpdateName] = useState('');
    const [updateId, setUpdateId] = useState(null);

    useEffect(() => {
        fetchNames();
    }, []);

    const fetchNames = async () => {
        try {
            const response = await axios.get('http://localhost:8080/names');
            setNames(response.data);
        } catch (error) {
            console.error("Error fetching names:", error);
        }
    };

    const createName = async () => {
        console.log("Creating name:", newName); // Debugging statement
        try {
            const response = await axios.post('http://localhost:8080/names', { name: newName });
            console.log("Response from server:", response); // Debugging statement
            setNewName('');
            fetchNames();
        } catch (error) {
            console.error("Error creating name:", error);
        }
    };

    const startUpdate = (id, currentName) => {
        setUpdateId(id);
        setUpdateName(currentName);
    };

    const saveUpdate = async () => {
        console.log("Updating name with ID:", updateId); // Debugging statement
        try {
            const response = await axios.put(`http://localhost:8080/names/${updateId}`, { name: updateName });
            console.log("Response from server:", response); // Debugging statement
            setUpdateId(null);
            setUpdateName('');
            fetchNames();
        } catch (error) {
            console.error("Error updating name:", error);
        }
    };


    const deleteName = async (id) => {
        try {
            await axios.delete(`http://localhost:8080/names/${id}`);
            fetchNames();
        } catch (error) {
            console.error("Error deleting name:", error);
        }
    };

    return (
        <div>
            <h1>Names List</h1>
            <input
                type="text"
                value={newName}
                onChange={(e) => setNewName(e.target.value)}
                placeholder="Add new name"
            />
            <button onClick={createName}>Create</button>
            <ul>
                {names && names.length > 0 ? (  // Check if names array has data
                    names.map((name) => (
                        <li key={name.id}>
                            {updateId === name.id ? (
                                <div>
                                    <input
                                        type="text"
                                        value={updateName}
                                        onChange={(e) => setUpdateName(e.target.value)}
                                    />
                                    <button onClick={saveUpdate}>Save</button>
                                </div>
                            ) : (
                                <div>
                                    {name.name}
                                    <button onClick={() => startUpdate(name.id, name.name)}>Update</button>
                                    <button onClick={() => deleteName(name.id)}>Delete</button>
                                </div>
                            )}
                        </li>
                    ))
                ) : (
                    <p>No names available</p>
                )}
            </ul>
        </div>
    );
}

export default App;
