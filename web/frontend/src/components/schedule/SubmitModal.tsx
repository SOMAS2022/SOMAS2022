import Backdrop from "@mui/material/Backdrop";
import Box from "@mui/material/Box";
import Modal from "@mui/material/Modal";
import Fade from "@mui/material/Fade";
import Typography from "@mui/material/Typography";

const style = {
    position: "absolute" as const,
    top: "50%",
    left: "50%",
    transform: "translate(-50%, -50%)",
    width: 400,
    bgcolor: "background.paper",
    border: "2px solid #000",
    boxShadow: 24,
    p: 4,
};

interface SubmitModalProps {
    status: number,
    text: string,
    initialOpen: boolean,
    toggle: () => void,
}

export default function SubmitModal({status, text, initialOpen, toggle}: SubmitModalProps) {
    return (
        <div>
            <Modal
                aria-labelledby="transition-modal-title"
                aria-describedby="transition-modal-description"
                open={initialOpen}
                onClose={toggle}
                closeAfterTransition
                BackdropComponent={Backdrop}
                BackdropProps={{
                    timeout: 500,
                }}
            >
                <Fade in={initialOpen}>
                    <Box sx={style}>
                        <Typography id="transition-modal-title" variant="h6" component="h2">
                            {status}
                        </Typography>
                        <Typography id="transition-modal-description" sx={{ mt: 2 }}>
                            {text}
                        </Typography>
                    </Box>
                </Fade>
            </Modal>
        </div>
    );
}