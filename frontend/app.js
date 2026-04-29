const API_URL = '/api';

document.addEventListener('DOMContentLoaded', () => {
    loadFeedback();
    
    document.getElementById('feedbackForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const feedback = {
            name: document.getElementById('name').value,
            email: document.getElementById('email').value,
            subject: document.getElementById('subject').value,
            message: document.getElementById('message').value
        };

        try {
            const response = await fetch(`${API_URL}/feedback`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(feedback)
            });

            if (response.ok) {
                document.getElementById('feedbackForm').reset();
                loadFeedback();
            } else {
                alert('Failed to submit feedback');
            }
        } catch (error) {
            console.error('Error:', error);
            alert('Error submitting feedback');
        }
    });
});

async function loadFeedback() {
    try {
        const response = await fetch(`${API_URL}/feedback`);
        const feedbacks = await response.json();
        displayFeedback(feedbacks);
    } catch (error) {
        console.error('Error loading feedback:', error);
    }
}

function displayFeedback(feedbacks) {
    const container = document.getElementById('feedbackList');
    container.innerHTML = '';

    if (!feedbacks || feedbacks.length === 0) {
        container.innerHTML = '<p>No feedback yet.</p>';
        return;
    }

    feedbacks.forEach(feedback => {
        const div = document.createElement('div');
        div.className = 'feedback-item';
        div.innerHTML = `
            <h3>${escapeHtml(feedback.name)} - ${escapeHtml(feedback.subject)}</h3>
            <p>Email: ${escapeHtml(feedback.email)}</p>
            <p class="message">${escapeHtml(feedback.message)}</p>
        `;
        container.appendChild(div);
    });
}

function escapeHtml(text) {
    const map = {
        '&': '&amp;',
        '<': '<',
        '>': '>',
        '"': '"',
        "'": '&#039;'
    };
    return text.replace(/[&<>"']/g, m => map[m]);
}